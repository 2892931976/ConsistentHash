package main


import (
    "hash/crc32"
    "sort"
)


type HashFunc func(date []byte) uint32

type HashRing []uint32

func(r HashRing) Len() int { return len(r) }
func(r HashRing) Less(i, j int) bool { return c[i] < c[j] }
func(r HahsRing) Swap(i, j int) { c[i], c[j] = c[j], c[i] }


type ConHash struct {
    hashfunc  HashFunc
    hashring  HashRing
    node_map_num      map[string] uint32  //physical node---> node num
    vnode_map_node    map[uint32] string  //vnode key ---->physical node
    
}



func ConHashNew() *ConHash {
    return &ConHash {
       hashfunc: crc32.ChecksumIEEE  
       hashring: HashRing{}
       node_map_num: make(map[string] uint32)
       vnode_map_node: make(map[uint32] string)
    }
}

func (ch *ConHash) NodeHash(data []byte) uint32 {
    return ch.hashfunc(data) 
}

//vnode format:  "n_name_%d"
func (ch *ConHash) NodeAdd(n_name string, vn_num uint32) {
    var name string
    var key uint32
    ch.node_map_num[n_name] = vn_num
    for i := 0; i < vn_num; i++ {
        name = n_name + "_" + string(i)
        fmt.Println(i, name)
        key = ch.NodeHash([]byte(name))

    }
}

func (ch *ConHash) NodeRemove(n_name string) {

}


func (ch *ConHash) NodeLookup (n_name string) {

}
