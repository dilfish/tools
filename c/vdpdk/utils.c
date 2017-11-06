/*
static const char*
arp_op_name(uint16_t arp_op) {
	switch(arp_op) {
		case ARP_OP_REQUEST:
			return "ARP REQUEST";
		case ARP_OP_REPLY:
			return "ARP REPLY";
		default:
			break;
	}
	return "Unkown ARP OP";
}


static const char*
ip_proto_name(uint16_t ip_proto) {
	switch(ip_proto) {
		case 4: // IPv4
			return "IPv4";
		case 6: // TCP
			return "TCP";
		case 17: // UDP
			return "UDP";
		default:
			break;
	}
	return "Unkown IP OP";
}
  
static void  
ipv4_addr_to_dot(uint32_t be_ipv4_addr, char *buf) { 
    uint32_t ipv4_addr;  
    ipv4_addr = rte_be_to_cpu_32(be_ipv4_addr);  
    sprintf(buf, "%d.%d.%d.%d", (ipv4_addr >> 24) & 0xFF,  
        (ipv4_addr >> 16) & 0xFF, (ipv4_addr >> 8) & 0xFF,  
        ipv4_addr & 0xFF);  
}
  
static void  
ether_addr_dump(const char *what, const struct ether_addr *ea) { 
    char buf[ETHER_ADDR_FMT_SIZE];  
    ether_format_addr(buf, ETHER_ADDR_FMT_SIZE, ea);  
    if (what)  
        printf("%s", what);  
    printf("%s", buf);  
}  
  
static void  
ipv4_addr_dump(const char *what, uint32_t be_ipv4_addr) { 
    char buf[16];  
    ipv4_addr_to_dot(be_ipv4_addr, buf);  
    if (what)  
        printf("%s", what);  
    printf("%s", buf);  
}  
*/
  
/*
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
*/
