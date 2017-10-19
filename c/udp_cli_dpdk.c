#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <stdlib.h>
#include <errno.h>
#include <sys/queue.h>

#include <rte_common.h>
#include <rte_memory.h>
#include <rte_memzone.h>
#include <rte_launch.h>
#include <rte_eal.h>
#include <rte_per_lcore.h>
#include <rte_ip.h>
#include <rte_ethdev.h>
#include <rte_ether.h>
#include <rte_lcore.h>
#include <rte_cycles.h>
#include <rte_timer.h>
#include <rte_udp.h>
#include <rte_debug.h>
#include <rte_arp.h>


#define MAX_PKT_BURST (64)


struct rte_mempool *glb_mp = NULL;
struct ether_addr eth_src_addr;
struct ether_addr eth_rt_addr;
// local addr 10.0.2.6
uint32_t ip_src_addr = 0xa000206;
// route addr 10.0.2.1
uint32_t ip_dst_addr = 0xa000201;
// libsm.com addr 119.28.76.79
uint32_t ip_libsm_addr = 0x771c4c4f;



static uint16_t  
ipv4_hdr_cksum(struct ipv4_hdr *ip_h)  
{  
    uint16_t *v16_h;  
    uint32_t ip_cksum;  
  
    /* 
     * Compute the sum of successive 16-bit words of the IPv4 header, 
     * skipping the checksum field of the header. 
     */  
    v16_h = (unaligned_uint16_t *) ip_h;  
    ip_cksum = v16_h[0] + v16_h[1] + v16_h[2] + v16_h[3] +  
        v16_h[4] + v16_h[6] + v16_h[7] + v16_h[8] + v16_h[9];  
  
    /* reduce 32 bit checksum to 16 bits and complement it */  
    ip_cksum = (ip_cksum & 0xffff) + (ip_cksum >> 16);  
    ip_cksum = (ip_cksum & 0xffff) + (ip_cksum >> 16);  
    ip_cksum = (~ip_cksum) & 0x0000FFFF;  
    return (ip_cksum == 0) ? 0xFFFF : (uint16_t) ip_cksum;  
}  


static int
send_arp_request(void) {
	printf("send arp request\n");
	struct rte_mbuf *m = NULL;
	struct ether_hdr *eth = NULL;
	struct arp_hdr *arp = NULL;
	int ret;
	m = rte_pktmbuf_alloc(glb_mp);
	if (m == NULL) {
		printf("alloc arp mbuf failed\n");
		return -1;
	}
	eth = rte_pktmbuf_mtod(m, struct ether_hdr*);
	arp = (struct arp_hdr*)&eth[1];
	memset(&eth->d_addr, 0xFF, 6);
	ether_addr_copy(&eth_src_addr, &eth->s_addr);
	eth->ether_type = htons(ETHER_TYPE_ARP);
	memset(arp, 0, sizeof(struct arp_hdr));
	rte_memcpy(arp->arp_data.arp_sha.addr_bytes, &eth_src_addr, 6);
	arp->arp_data.arp_sip = htonl(ip_src_addr);
	arp->arp_data.arp_tip = htonl(ip_dst_addr);
	memset(arp->arp_data.arp_tha.addr_bytes, 0, 6);
	const uint16_t ETH_HW_TYPE = 1;
	arp->arp_hrd = htons(ETH_HW_TYPE);
	arp->arp_pro = htons(ETHER_TYPE_IPv4);
	arp->arp_hln = 6;
	arp->arp_pln = 4;
	arp->arp_op = htons(ARP_OP_REQUEST);
	uint32_t sz = sizeof(struct arp_hdr) + sizeof(struct ether_hdr);
	m->pkt_len = sz;
	m->data_len = sz;
	ret = rte_eth_tx_burst(0, 0, &m, 1);
	if (ret < 0) {
		printf("send arp request err\n");
		rte_pktmbuf_free(m);
	}
	return ret;
}

static int
recv_arp_response(struct rte_mbuf **mbufs, uint16_t nb_rx) {
	int i;
	struct rte_mbuf *m;
	struct arp_hdr *arp;
	for (i = 0;i < nb_rx; i ++) {
		m = mbufs[i];
		struct ether_hdr *eth = NULL;
		eth = rte_pktmbuf_mtod(m, struct ether_hdr*);
		if (eth->ether_type != htons(ETHER_TYPE_ARP)) {
			continue;
		}
		arp = (struct arp_hdr*)&eth[1];
		if (arp->arp_op != htons(ARP_OP_REPLY)) {
			continue;
		}
		printf("recved arp response, pktlen %u\n", m->pkt_len);
		printf("ip src %u, dst %u\n", arp->arp_data.arp_sip, arp->arp_data.arp_tip);
		printf("mac src:");
		for (i = 0;i < 6;i ++) {
			printf("%x:", arp->arp_data.arp_sha.addr_bytes[i]);
		}
		ether_addr_copy(&arp->arp_data.arp_sha, &eth_rt_addr);
		printf(", mac dst:");
		for (i = 0;i < 6;i ++) {
			printf("%x:", arp->arp_data.arp_tha.addr_bytes[i]);
		}
		printf("\n");
		return 0;
	}
	return -1;
}


static int
send_udp(void) {
	printf("start to send udp\n");
	struct rte_mbuf *m = NULL;
	struct ether_hdr *eth = NULL;
	int i, ret;
	struct ipv4_hdr *ip_h;
	struct udp_hdr *u_hdr;
	m = rte_pktmbuf_alloc(glb_mp);
	if (m == NULL) {
		printf("alloc arp mbuf failed\n");
		return -1;
	}
	eth = rte_pktmbuf_mtod(m, struct ether_hdr*);
	ether_addr_copy(&eth_rt_addr, &eth->d_addr);
	ether_addr_copy(&eth_src_addr, &eth->s_addr);
	eth->ether_type = htons(ETHER_TYPE_IPv4);
	ip_h = (struct ipv4_hdr*)&eth[1];
	ip_h->hdr_checksum = 0;
	ip_h->src_addr = htonl(ip_src_addr);
	ip_h->dst_addr = htonl(ip_libsm_addr);
	ip_h->version_ihl = 69;
	ip_h->type_of_service = 0;
	// 36 == abcdefg\0
	// 20 + 8 + 8
	ip_h->total_length = htons(8 + sizeof(struct udp_hdr) + sizeof(struct ipv4_hdr));
	ip_h->packet_id = (uint16_t)random();
	ip_h->time_to_live = 64;
	// flag set, no frag
	ip_h->fragment_offset = 64;
	ip_h->next_proto_id = 17;
	ip_h->hdr_checksum = ipv4_hdr_cksum(ip_h);
	u_hdr = (struct udp_hdr*)&ip_h[1];
	u_hdr->src_port = htons(9999);
	u_hdr->dst_port = htons(9999);
	u_hdr->dgram_cksum = 0;
	u_hdr->dgram_len = htons(8 + sizeof(struct udp_hdr));
	char *msg = (char*)u_hdr + sizeof(struct udp_hdr);
	for (i = 0;i < 7;i ++) {
		msg[i] = 'a' + i;
	}
	msg[7] = '\n';
	uint32_t sz = sizeof(struct ether_hdr) + sizeof(struct ipv4_hdr) + sizeof(struct udp_hdr);
	sz = sz + 8;
	m->pkt_len = sz;
	m->data_len = sz;
	ret = rte_eth_tx_burst(0, 0, &m, 1);
	if (ret < 0) {
		rte_pktmbuf_free(m);
	}
	return ret;
}


static int
recv_udp_response(struct rte_mbuf **mbufs, uint16_t nb_rx) {
	uint16_t i;
	struct ether_hdr *eth_h;
	struct ipv4_hdr *ip_h;
	uint16_t eth_type;
	for (i = 0;i < nb_rx;i ++) {
		struct rte_mbuf *m = mbufs[i];
		eth_h = rte_pktmbuf_mtod(m, struct ether_hdr*);
		eth_type = rte_be_to_cpu_16(eth_h->ether_type);
		uint16_t l2_len;
		l2_len = sizeof(struct ether_hdr);
		if (eth_type != ETHER_TYPE_IPv4) {
			continue;
		}
		ip_h = (struct ipv4_hdr*)((char*)eth_h + l2_len);
		if (ip_h->next_proto_id != 17) {
			continue;
		}
		if (ip_h->dst_addr != htonl(ip_src_addr)) {
			continue;
		}
		printf("ipaddr is src 0x%x, dst 0x%x\n", ip_h->src_addr, ip_h->dst_addr);
		struct udp_hdr *u_hdr;
		u_hdr = (struct udp_hdr*)((rte_pktmbuf_mtod(m, char*) + l2_len + sizeof(struct ipv4_hdr)));
		char *msg = (char*)u_hdr + sizeof(struct udp_hdr);
		printf("msg header is %c, %c, %c\n", msg[0], msg[1], msg[2]);
		return 0;
	}
	return -1;
}


#define TIMER_RESOLUTION_CYCLES 20000000ULL /* around 10ms at 2 Ghz */

static struct rte_timer timer0;
int status = 0;

/* timer0 callback */
static void
timer0_cb(__attribute__((unused)) struct rte_timer *tim,
	  __attribute__((unused)) void *arg)
{
	unsigned lcore_id = rte_lcore_id();
	int ret;
	struct rte_mbuf *pkts_burst[MAX_PKT_BURST];
	int16_t i, nb_rx;

	printf("%s() on lcore %u\n", __func__, lcore_id);
	// port, queue
    	nb_rx = rte_eth_rx_burst(0, 0, pkts_burst, MAX_PKT_BURST);
	if (nb_rx < 0) {
		printf("rx error\n");
		rte_timer_stop(tim);
		return;
	}
	printf("nb_rx is %d\n", nb_rx);
	for (i = 0;i < nb_rx;i ++) {
		struct ether_hdr *eth_hdr = rte_pktmbuf_mtod(pkts_burst[i], struct ether_hdr*);
		printf("eth type is %u\n", htons(eth_hdr->ether_type));
		if (htons(eth_hdr->ether_type) == 2048) {
			struct ipv4_hdr *ip_h = (struct ipv4_hdr*)&eth_hdr[1];
			printf("ipv4 next type %u\n", ip_h->next_proto_id);
			printf("ipv4 src dst 0x%x, 0x%x\n", ip_h->src_addr, ip_h->dst_addr);
		}
	}
	if (status == 0) {
		send_arp_request();
	}
	ret = recv_arp_response(pkts_burst, nb_rx);
	if (ret == 0) {
		status = 1;
	}
	if (status == 1) {
		ret = send_udp();
		if (ret < 0) {
			printf("send udp error\n");
			rte_timer_stop(tim);
			return;
		}
	}
	ret = recv_udp_response(pkts_burst, nb_rx);
	if (ret == 0) {
		printf("recv udp response\n");
		rte_timer_stop(tim);
		return;
	}
}

static __attribute__((noreturn)) int
lcore_mainloop(__attribute__((unused)) void *arg)
{
	uint64_t prev_tsc = 0, cur_tsc, diff_tsc;
	unsigned lcore_id;

	lcore_id = rte_lcore_id();
	printf("Starting mainloop on core %u\n", lcore_id);

	while (1) {
		cur_tsc = rte_rdtsc();
		diff_tsc = cur_tsc - prev_tsc;
		if (diff_tsc > TIMER_RESOLUTION_CYCLES) {
			rte_timer_manage();
			prev_tsc = cur_tsc;
		}
	}
}


static const struct rte_eth_conf port_conf_default = {  
    .rxmode = { .max_rx_pkt_len = ETHER_MAX_LEN }  
}; 


#define RX_RING_SIZE 128  
#define TX_RING_SIZE 512


static inline int  
port_init(uint8_t port, struct rte_mempool *mbuf_pool)  
{  
    struct rte_eth_conf port_conf = port_conf_default;  
    const uint16_t rx_rings = 1, tx_rings = 1;  
    int retval;  
    uint16_t q;  
  
    if (port >= rte_eth_dev_count())  
        return -1;  
  
    /* Configure the Ethernet device. */  
    retval = rte_eth_dev_configure(port, rx_rings, tx_rings, &port_conf);  
    if (retval != 0)  
        return retval;  
  
    /* Allocate and set up 1 RX queue per Ethernet port. */  
    for (q = 0; q < rx_rings; q++) {  
        retval = rte_eth_rx_queue_setup(port, q, RX_RING_SIZE,  
                rte_eth_dev_socket_id(port), NULL, mbuf_pool);  
        if (retval < 0)  
            return retval;  
    }  
  
    /* Allocate and set up 1 TX queue per Ethernet port. */  
    for (q = 0; q < tx_rings; q++) {  
        retval = rte_eth_tx_queue_setup(port, q, TX_RING_SIZE,  
                rte_eth_dev_socket_id(port), NULL);  
        if (retval < 0)  
            return retval;  
    }  
  
    /* Start the Ethernet port. */  
    retval = rte_eth_dev_start(port);  
    if (retval < 0)  
        return retval;  
  
    /* Display the port MAC address. */  
    rte_eth_macaddr_get(port, &eth_src_addr);  
    printf("Port %u MAC: %02" PRIx8 " %02" PRIx8 " %02" PRIx8  
               " %02" PRIx8 " %02" PRIx8 " %02" PRIx8 "\n",  
            (unsigned)port,  
            eth_src_addr.addr_bytes[0], eth_src_addr.addr_bytes[1],  
            eth_src_addr.addr_bytes[2], eth_src_addr.addr_bytes[3],  
            eth_src_addr.addr_bytes[4], eth_src_addr.addr_bytes[5]);  
  
    /* Enable RX in promiscuous mode for the Ethernet device. */  
    rte_eth_promiscuous_enable(port);  
  
    return 0;  
}  

int
main(int argc, char **argv)
{
	int ret;
	uint64_t hz;
	unsigned lcore_id;
	unsigned nb_port;
	ret = rte_eal_init(argc, argv);
	if (ret < 0)
		rte_panic("Cannot init EAL\n");
	rte_timer_subsystem_init();
	rte_timer_init(&timer0);
	printf("sizeof arp is %lu\n", sizeof(struct arp_hdr) + sizeof(struct ether_hdr));
	nb_port = rte_eth_dev_count();
	if (nb_port <= 0) {
		rte_panic("get eth dev error\n");
	}
	glb_mp = rte_pktmbuf_pool_create("pktpool", 8192*nb_port, 250, 0, RTE_MBUF_DEFAULT_BUF_SIZE, rte_socket_id());
	if (glb_mp == NULL) {
		rte_panic("cannot create mem pool\n");
	}
	if (port_init(0, glb_mp) != 0) {
		rte_panic("cannot port init\n");
	}
	hz = rte_get_timer_hz();
	lcore_id = rte_lcore_id();
	rte_timer_reset(&timer0, hz*2, PERIODICAL, lcore_id, timer0_cb, NULL);
	(void) lcore_mainloop(NULL);
	return 0;
}
