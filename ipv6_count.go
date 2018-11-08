package tools

import (
    "strings"
    "fmt"
    "strconv"
    "sort"
)


type StateCount struct {
    Name string
    Count uint64
}


type StateCountSlice []StateCount
func (scs StateCountSlice) Len() int {
    return len(scs)
}


func (scs StateCountSlice) Less (i, j int) bool {
    return scs[i].Count > scs[j].Count
}


func (scs StateCountSlice) Swap (i, j int) {
    scs[i], scs[j] = scs[j], scs[i]
}


type IPv6Counter struct {
    stateList []string
    current string
    stateCountMap map[string]uint64
    sortCountList []StateCount
}

func NewIPv6Counter() *IPv6Counter {
    tmp := []string {
        "AD","AE","AF","AG","AI","AL","AM","AO","AR","AS","AT",
        "AU","AW","AX","AZ","BA","BB","BD","BE","BF","BG","BH",
        "BI","BJ","BL","BM","BN","BO","BQ","BR","BS","BT","BW",
        "BY","BZ","CA","CD","CF","CG","CH","CI","CK","CL","CM",
        "CN","CO","CR","CU","CV","CW","CY","CZ","DE","DJ","DK",
        "DM","DO","DZ","EC","EE","EG","ER","ES","ET","FI","FJ",
        "FK","FM","FO","FR","GA","GB","GD","GE","GF","GG","GH",
        "GI","GL","GM","GN","GP","GQ","GR","GT","GU","GW","GY",
        "HK","HN","HR","HT","HU","ID","IE","IL","IM","IN","IO",
        "IQ","IR","IS","IT","JE","JM","JO","JP","KE","KG","KH",
        "KI","KM","KN","KP","KR","KW","KY","KZ","LA","LB","LC",
        "LI","LK","LR","LS","LT","LU","LV","LY","MA","MC","MD",
        "ME","MF","MG","MH","MK","ML","MM","MN","MO","MP","MQ",
        "MR","MT","MU","MV","MW","MX","MY","MZ","NA","NC","NE",
        "NF","NG","NI","NL","NO","NP","NR","NU","NZ","OM","PA",
        "PE","PF","PG","PH","PK","PL","PM","PR","PS","PT","PW",
        "PY","QA","RE","RO","RS","RU","RW","SA","SB","SC","SD",
        "SE","SG","SI","SK","SL","SM","SN","SO","SR","SS","ST",
        "SV","SX","SY","SZ","TC","TD","TG","TH","TJ","TK","TL",
        "TM","TN","TO","TR","TT","TV","TW","TZ","UA","UG","US",
        "UY","UZ","VA","VC","VE","VG","VI","VN","VU","WF","WS",
        "YE","YT","ZA","ZM","ZW",
    }
    var ipv6c IPv6Counter
    for _, t := range tmp {
        l := strings.ToLower(t)
        ipv6c.stateList = append(ipv6c.stateList, l)
    }
    ipv6c.stateCountMap = make(map[string]uint64)
    return &ipv6c
}


func (ipv6c *IPv6Counter) count (line string) error {
    if len(line) >= 1 && line[0] == '#' {
        return nil
    }
    arr := strings.Split(line, "/")
    if len(arr) != 2 {
        fmt.Println("bad format", arr)
        return ErrBadFmt
    }
    mask, err := strconv.ParseUint(arr[1], 10, 32)
    if err != nil {
        return err
    }
    mask = 64 - mask
    total := ipv6c.stateCountMap[ipv6c.current]
    total = total + (1 << mask)
    ipv6c.stateCountMap[ipv6c.current] = total
    return nil
}


func (ipv6c *IPv6Counter) get_files() error {
    base := "http://ipverse.net/ipblocks/data/countries/"
    for _, state := range ipv6c.stateList {
        ipv6c.current = state
        uri := base + state + "-ipv6.zone"
        err := GetLine(uri, ipv6c.count)
        if err != nil {
            return err
        }
    }
    return nil
}


func (ipv6c *IPv6Counter) sort_counter () {
    scList := make([]StateCount, 0)
    for k, v := range ipv6c.stateCountMap {
        var sc StateCount
        sc.Name = k
        sc.Count = v
        scList = append(scList, sc)
    }
    sort.Sort(StateCountSlice(scList))
    ipv6c.sortCountList = scList
}


func (ipv6c *IPv6Counter) Renew() error {
    ipv6c.current = ""
    ipv6c.sortCountList = nil
    ipv6c.stateCountMap = make(map[string]uint64)
    err := ipv6c.get_files()
    if err != nil {
        return err
    }
    ipv6c.sort_counter()
    return nil
}


func (ipv6c *IPv6Counter) String() string {
    str := ""
    for _, state := range ipv6c.sortCountList {
        str = str + state.Name + "\t" + strconv.FormatUint(state.Count, 10) + "\n"
    }
    return str
}
