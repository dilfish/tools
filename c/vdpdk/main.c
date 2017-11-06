#include "header.h"


struct ether_addr eth_addrs[2]; // 0==rx, 1==tx
struct rte_mempool *global_mbuf_pool = NULL;
static const struct rte_eth_conf port_conf_default = {  
    .rxmode = { .max_rx_pkt_len = ETHER_MAX_LEN }  
};  


struct ether_addr eth_src_addr;
struct ether_addr eth_rt_addr;


#include "utils.c"


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
handle_ipv4(struct rte_mbuf *pkt, struct ether_hdr *eth_h) {
	struct ipv4_hdr *ip_h = (struct ipv4_hdr*)((char*)eth_h + pkt->l2_len);
	struct icmp_hdr *icmp_h = (struct icmp_hdr*) ((char*)ip_h + sizeof(struct ipv4_hdr));
	if (ip_h->next_proto_id == IPPROTO_ICMP) {
		return handle_icmp(ip_h, icmp_h);
	}
	// check ip
	// UDP_PROTO_ID = 17
	return 0;
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
	if (rte_be_to_cpu_16(arp_h->arp_op) != ARP_OP_REQUEST) {
		return -2;
	}
	// check ip
	ether_addr_copy(&eth_h->s_addr, &eth_h->d_addr);
	ether_addr_copy(&eth_addrs[0], &eth_h->s_addr);
	arp_h->arp_op = rte_cpu_to_be_16(ARP_OP_REPLY);
	ether_addr_copy(&arp_h->arp_data.arp_tha, &eth_addrs[0]);
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
				ret = handle_ipv4(mbufs[i], eth_h); 
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
port_init(uint8_t portid, struct rte_mempool *mbuf_pool, struct ether_addr *addr) {
	struct rte_eth_conf port_conf = port_conf_default;
	int ret;
	ret = rte_eth_dev_configure(portid, 1, 1, &port_conf);
	if (ret != 0) {
		return ret;
	}
	uint16_t socket = rte_eth_dev_socket_id(portid);
	ret = rte_eth_rx_queue_setup(portid, 0, RX_RING_SIZE, socket, NULL, mbuf_pool);
	if (ret != 0) {
		return ret;
	}
	ret = rte_eth_tx_queue_setup(portid, 0, TX_RING_SIZE, socket, NULL);
	if (ret != 0) {
		return ret;
	}
	ret = rte_eth_dev_start(portid);
	if (ret != 0) {
		return ret;
	}
    rte_eth_macaddr_get(portid, addr);  
	rte_eth_promiscuous_enable(portid);
    return 0;  
}  


int
main(int argc, char **argv) {
	int i, ret, nb_ports;
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
	for (i = 0;i < 2;i ++) {
		ret = port_init(i, global_mbuf_pool, &eth_addrs[i]);
		if (ret < 0) {
			rte_panic("init port error\n");
		}
	}
	handle_ports();
	return 0;
}
