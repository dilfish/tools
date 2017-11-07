#ifndef HEADER_H
#define HEADER_H

// standard c library
#include <stdarg.h>  
#include <string.h>  
#include <stdio.h>  
#include <errno.h>  
#include <stdint.h>  
#include <inttypes.h>  
#include <stdlib.h>
// linux library  
#include <unistd.h>
#include <sys/queue.h>  
#include <sys/stat.h>  
// dpdk library 
#include <rte_eal.h>
#include <rte_per_lcore.h>
#include <rte_ip.h>
#include <rte_cycles.h>
#include <rte_timer.h>
#include <rte_debug.h>
#include <rte_common.h>  
#include <rte_byteorder.h>  
#include <rte_log.h>  
#include <rte_debug.h>  
#include <rte_cycles.h>  
#include <rte_launch.h>
#include <rte_per_lcore.h>  
#include <rte_lcore.h>  
#include <rte_atomic.h>  
#include <rte_branch_prediction.h>  
#include <rte_ring.h>  
#include <rte_memory.h>  
#include <rte_memzone.h>
#include <rte_udp.h>
#include <rte_mempool.h>  
#include <rte_mbuf.h>  
#include <rte_ether.h>  
#include <rte_ethdev.h>  
#include <rte_arp.h>  
#include <rte_ip.h>  
#include <rte_icmp.h>  
#include <rte_string_fns.h>  


#define is_multicast_ipv4_addr(ipv4_addr) \
    (((rte_be_to_cpu_32((ipv4_addr)) >> 24) & 0x000000FF) == 0xE0)  
  
#define MAX_PKT_BURST 512  
#define RX_RING_SIZE 128  
#define TX_RING_SIZE 512  
  
#define NUM_MBUFS 8191  

#endif
