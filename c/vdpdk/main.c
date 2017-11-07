// virtual box test for lib dpdk
// virtual box has two network mode:
// 1. bridge: host could send and recv with guest
// 2. nat: guest could send and recv with remote
// so we use two ports combined

// 0. use nat port as port 0 and bridge as port 1
// 1. detect route mac of port 0
// 2. store route mac of port 0
// 3. listen on port 1
// 4. host send udp abcdefg\n string to guest via port 1
// 5. port 1 received this string and send it
// to remote addr via port 0

#include "header.h"


struct rte_mempool *global_mbuf_pool = NULL;
static const struct rte_eth_conf port_conf_default = {  
    .rxmode = { .max_rx_pkt_len = ETHER_MAX_LEN }  
};  


int status = 0;


struct ether_addr eth_tx_addr;
struct ether_addr eth_rx_addr;
struct ether_addr eth_rt_addr;

// 10.0.2.1 route ip for port 0
uint32_t ip_rt_addr = 0xa000201;
// send ip 10.0.2.17, port 0
uint32_t ip_tx_addr = 0xa000211;
// recv ip 100.100.33.196, port 1
uint32_t ip_rx_addr = 0x646421c4;
// 119.28.76.79 remote server
uint32_t ip_libsm_addr = 0x771c4c4f;


static uint16_t
ipv4_hdr_cksum(struct ipv4_hdr *ip_h) {
    uint16_t *v16_h;
    uint32_t ip_cksum;
    v16_h = (unaligned_uint16_t *) ip_h;
    ip_cksum = v16_h[0] + v16_h[1] + v16_h[2] + v16_h[3] +
        v16_h[4] + v16_h[6] + v16_h[7] + v16_h[8] + v16_h[9];
    ip_cksum = (ip_cksum & 0xffff) + (ip_cksum >> 16);
    ip_cksum = (ip_cksum & 0xffff) + (ip_cksum >> 16);
    ip_cksum = (~ip_cksum) & 0x0000FFFF;
    return (ip_cksum == 0) ? 0xFFFF : (uint16_t) ip_cksum;
}


static void
print_eth_addr(const unsigned char *addr) {
	int i;
	for (i = 0;i < 6;i ++) {
		printf("%x:", addr[i]);
	}
}

static int
send_arp_request(uint16_t portid) {
	printf("send arp request\n");
	struct rte_mbuf *m = NULL;
	struct ether_hdr *eth = NULL;
	struct arp_hdr *arp = NULL;
	int ret;
	m = rte_pktmbuf_alloc(global_mbuf_pool);
	if (m == NULL) {
		printf("alloc arp mbuf failed\n");
		return -1;
	}
	eth = rte_pktmbuf_mtod(m, struct ether_hdr*);
	arp = (struct arp_hdr*)&eth[1];
	memset(&eth->d_addr, 0xFF, 6);
	ether_addr_copy(&eth_tx_addr, &eth->s_addr);
	eth->ether_type = htons(ETHER_TYPE_ARP);
	memset(arp, 0, sizeof(struct arp_hdr));
	rte_memcpy(arp->arp_data.arp_sha.addr_bytes, &eth_tx_addr, 6);
	arp->arp_data.arp_sip = htonl(ip_tx_addr);
	arp->arp_data.arp_tip = htonl(ip_rt_addr);
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
	ret = rte_eth_tx_burst(portid, 0, &m, 1);
	if (ret < 0) {
		printf("send arp request err\n");
		rte_pktmbuf_free(m);
	}
	return ret;
}


static int
send_udp(const char *buff, int buff_len) {
	printf("start to send udp\n");
	struct rte_mbuf *m = NULL;
	struct ether_hdr *eth = NULL;
	int ret;
	struct ipv4_hdr *ip_h;
	struct udp_hdr *u_hdr;
	m = rte_pktmbuf_alloc(global_mbuf_pool);
	if (m == NULL) {
		printf("alloc arp mbuf failed\n");
		return -1;
	}
	eth = rte_pktmbuf_mtod(m, struct ether_hdr*);
	ether_addr_copy(&eth_rt_addr, &eth->d_addr);
	ether_addr_copy(&eth_tx_addr, &eth->s_addr);
	eth->ether_type = htons(ETHER_TYPE_IPv4);
	ip_h = (struct ipv4_hdr*)&eth[1];
	ip_h->hdr_checksum = 0;
	ip_h->src_addr = htonl(ip_tx_addr);
	ip_h->dst_addr = htonl(ip_libsm_addr);
	ip_h->version_ihl = 69;
	ip_h->type_of_service = 0;
	// 36 == abcdefg\0
	// 20 + 8 + 8
	ip_h->total_length = htons(buff_len + sizeof(struct udp_hdr) + sizeof(struct ipv4_hdr));
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
	u_hdr->dgram_len = htons(buff_len + sizeof(struct udp_hdr));
	char *msg = (char*)u_hdr + sizeof(struct udp_hdr);
	memcpy(msg, buff, buff_len);
	uint32_t sz = sizeof(struct ether_hdr) + sizeof(struct ipv4_hdr) + sizeof(struct udp_hdr);
	sz = sz + buff_len;
	m->pkt_len = sz;
	m->data_len = sz;
	ret = rte_eth_tx_burst(0, 0, &m, 1);
	if (ret < 0) {
		rte_pktmbuf_free(m);
	}
	return ret;
}


static int
handle_udp(struct ipv4_hdr *ip, uint16_t portid) {
	if (ip->dst_addr != htonl(ip_tx_addr) && ip->dst_addr != htonl(ip_rx_addr)) {
		// printf("not local ip 0x%x, 0x%x, 0x%x\n", htonl(ip->dst_addr), ip_tx_addr, ip_rx_addr);
		return -1;
	}
	printf("ipaddr src 0x%x, dst 0x%x\n", ip->src_addr, ip->dst_addr);
	struct udp_hdr *u_hdr;
	u_hdr = (struct udp_hdr*)((char*)ip + sizeof(struct ipv4_hdr));
	char *msg = (char*)u_hdr + sizeof(struct udp_hdr);
	printf("msg header is %c, %c, %c\n", msg[0], msg[1], msg[2]);
	printf("port id is %u\n", portid);
	if (portid != 0) {
		return send_udp(msg, 8);
	}
	status = 2;
	return -1;
}




static int
handle_icmp(struct ipv4_hdr *ip, struct icmp_hdr *icmp) {
	uint32_t cksum;
	if (icmp->icmp_type != IP_ICMP_ECHO_REQUEST && icmp->icmp_code != 0) {
		return -1;
	}
	if (is_multicast_ipv4_addr(ip->dst_addr)) {
		uint32_t ip_src, ip_addr = 0;
		ip_src = rte_be_to_cpu_32(ip_addr);
		if ((ip_src & 0x3) == 1) {
			ip_src = (ip_src & 0xFFFFFFFC) | 0x1;
		}
	}
	icmp->icmp_type = IP_ICMP_ECHO_REPLY;
	cksum = ~icmp->icmp_cksum & 0xffff;
	cksum += ~htons(IP_ICMP_ECHO_REQUEST << 8) & 0xffff;
	cksum += htons(IP_ICMP_ECHO_REPLY << 8);
	cksum = (cksum & 0xffff) + (cksum >> 16);
	cksum = (cksum & 0xffff) + (cksum >> 16);
	icmp->icmp_cksum = ~cksum;
	return 0;
}


static int
handle_ipv4(struct rte_mbuf *pkt, struct ether_hdr *eth_h, uint16_t portid) {
	struct ipv4_hdr *ip_h = (struct ipv4_hdr*)((char*)eth_h + pkt->l2_len);
	struct icmp_hdr *icmp_h = (struct icmp_hdr*) ((char*)ip_h + sizeof(struct ipv4_hdr));
	if (ip_h->next_proto_id == IPPROTO_ICMP) {
		return handle_icmp(ip_h, icmp_h);
	}
	if (ip_h->next_proto_id == 17) {
		return handle_udp(ip_h, portid);
	}
	return 0;
}


static int
handle_arp_resp(struct rte_mbuf *pkt, struct ether_hdr *eth_h) {
	struct arp_hdr *arp = (struct arp_hdr*)((char*)eth_h + pkt->l2_len);
	printf("arp response: ip src 0x%x, dst 0x%x\n", arp->arp_data.arp_sip, arp->arp_data.arp_tip);
	printf("mac src: ");
	print_eth_addr(arp->arp_data.arp_sha.addr_bytes);
	printf("mac dst:");
	print_eth_addr(arp->arp_data.arp_tha.addr_bytes);
	ether_addr_copy(&arp->arp_data.arp_sha, &eth_rt_addr);
	status = 1;
	// free data
	return -1;
}


static int
handle_arp(struct rte_mbuf *pkt, struct ether_hdr *eth_h) {
	struct arp_hdr *arp_h = NULL;
	uint32_t ip_addr;
	arp_h = (struct arp_hdr*)((char*)eth_h + pkt->l2_len);
	if (rte_be_to_cpu_16(arp_h->arp_hrd) != ARP_HRD_ETHER ||
		rte_be_to_cpu_16(arp_h->arp_pro) != ETHER_TYPE_IPv4 ||
		arp_h->arp_hln != 6 || arp_h->arp_pln != 4) {
		return -1;
	}
	if (rte_be_to_cpu_16(arp_h->arp_op) == ARP_OP_REPLY) {
		return handle_arp_resp(pkt, eth_h);
	}
	if (rte_be_to_cpu_16(arp_h->arp_op) != ARP_OP_REQUEST) {
		return -2;
	}
	// check ip
	ether_addr_copy(&eth_h->s_addr, &eth_h->d_addr);
	ether_addr_copy(&eth_rx_addr, &eth_h->s_addr);
	arp_h->arp_op = rte_cpu_to_be_16(ARP_OP_REPLY);
	ether_addr_copy(&arp_h->arp_data.arp_tha, &eth_rx_addr);
    ether_addr_copy(&arp_h->arp_data.arp_sha, &arp_h->arp_data.arp_tha);  
    ether_addr_copy(&eth_h->s_addr, &arp_h->arp_data.arp_sha);  
    ip_addr = arp_h->arp_data.arp_sip;  
    arp_h->arp_data.arp_sip = arp_h->arp_data.arp_tip;  
    arp_h->arp_data.arp_tip = ip_addr;  
	return 0;
}


static uint16_t 
check_vlan(struct rte_mbuf *pkt) {
	struct ether_hdr *eth_h = NULL;
	struct vlan_hdr *vlan_h = NULL;
	uint16_t ether_type = 0;
	if (pkt == NULL) {
		return 0;
	}
	pkt->l2_len = sizeof(struct ether_hdr);
	eth_h = rte_pktmbuf_mtod(pkt, struct ether_hdr*);
	ether_type = rte_be_to_cpu_16(eth_h->ether_type);
	if (ether_type == ETHER_TYPE_VLAN) {
		pkt->l2_len += sizeof(struct vlan_hdr);
        vlan_h = (struct vlan_hdr *) 
                ((char *)eth_h + sizeof(struct ether_hdr));  
		ether_type = rte_be_to_cpu_16(vlan_h->eth_proto); 
	}
	return ether_type;
}


static int 
handle_loop(uint8_t portid) {
	struct rte_mbuf *mbufs[MAX_PKT_BURST], *send_back_mbufs[MAX_PKT_BURST];
	int ret, nb_rx = 0, i, send_backs = 0;
	uint16_t ether_type;
	nb_rx = rte_eth_rx_burst(portid, 0, mbufs, MAX_PKT_BURST);
	if (nb_rx <= 0) {
		return nb_rx;
	}
	for (i = 0;i < nb_rx; i ++) {
		ether_type = check_vlan(mbufs[i]);
		struct ether_hdr *eth_h = rte_pktmbuf_mtod(mbufs[i], struct ether_hdr *);
		switch(ether_type) {
			case ETHER_TYPE_ARP:
				ret = handle_arp(mbufs[i], eth_h); 
				break;
			case ETHER_TYPE_IPv4:
				ret = handle_ipv4(mbufs[i], eth_h, portid); 
				break;
			default:
				ret = -1;
				break;
		}
		if (ret == 0) {
			send_back_mbufs[send_backs] = mbufs[i];
			send_backs++;
		} else {
			rte_pktmbuf_free(mbufs[i]);
		}
	}
	return rte_eth_tx_burst(portid, 0, send_back_mbufs, send_backs);
}


static void
handle_ports(void) {
	while (1) {
		usleep(100*1000);
		handle_loop(0);
		handle_loop(1);
	}
}


static int
port_init(uint16_t portid, struct ether_addr *addr) {
	struct rte_eth_conf port_conf = port_conf_default;
	int ret, socket;
	socket = 0;
	ret = rte_eth_dev_configure(portid, 1, 1, &port_conf);
	if (ret != 0) {
		printf("dev configure error\n");
		return ret;
	}
	printf("port id %u\n", portid);
	ret = rte_eth_rx_queue_setup(portid, 0, RX_RING_SIZE, socket, NULL, global_mbuf_pool);
	if (ret != 0) {
		printf("rx queue error\n");
		return ret;
	}
	ret = rte_eth_tx_queue_setup(portid, 0, TX_RING_SIZE, socket, NULL);
	if (ret != 0) {
		printf("tx queue error\n");
		return ret;
	}
	ret = rte_eth_dev_start(portid);
	if (ret != 0) {
		printf("dev start error\n");
		return ret;
	}
    rte_eth_macaddr_get(portid, addr);  
	rte_eth_promiscuous_enable(portid);
    return 0;  
}  


static int
detect_mac(uint16_t port_id) {
	int ret;
	while (1) {
		ret = send_arp_request(port_id);
		if (ret < 0) {
			return -1;
		}
		usleep(1500*1000);
		handle_loop(0);
		if (status == 1) {
			printf("recved arp response\n");
			break;
		}
	}
	return 0;
	// after detect mac, we could test the avaliability of udp
	while (1) {
		ret = send_udp("abcdefg\n", 8);
		if (ret < 0) {
			printf("send udp error\n");
			return -1;
		}
		usleep(1500 * 1000);
		handle_loop(0);
		if (status == 2) {
			printf("rev udp resp\n");
			break;
		}
	}
	return 0;
}


int
main(int argc, char **argv) {
	int ret, nb_ports;
	uint16_t i;
	ret = rte_eal_init(argc, argv);
	if (ret < 0) {
		rte_panic("can not init eal\n");
	}
	nb_ports = rte_eth_dev_count();
	if (nb_ports != 2) {
		rte_panic("this demo could only works under 2 NICs\n");
	}
	global_mbuf_pool = rte_pktmbuf_pool_create("mbuf_pool", NUM_MBUFS * nb_ports,
		250, 0, RTE_MBUF_DEFAULT_BUF_SIZE, rte_socket_id());
	if (global_mbuf_pool == NULL) {
		rte_panic("could not alloc memory for pool\n");
	}
	struct ether_addr *ptr[2];
	ptr[0] = &eth_rx_addr;
	ptr[1] = &eth_tx_addr;
	for (i = 0;i < 2;i ++) {
		printf("init port %u\n", i);
		ret = port_init(i, ptr[i]);
		if (ret < 0) {
			rte_panic("init port error\n");
		}
	}
	ret = detect_mac(0);
	printf("detect_mac ret %d\n", ret);
	if (ret < 0) {
		printf("detect mac error\n");
		return -1;
	}
	handle_ports();
	return 0;
}
