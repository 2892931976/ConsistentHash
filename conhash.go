package main


import (
    "fmt"
    "hash/crc32"
    "strconv"
    "sort"
)

type HashFunc func(date []byte) uint32
type HashRing []uint32

func(r HashRing) Len() int { return len(r) }
func(r HashRing) Less(i, j int) bool { return r[i] < r[j] }
func(r HashRing) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type ConHash struct {
    hashfunc  HashFunc
    hashring  HashRing
    vnodes    int
    node_map_num      map[string] int  //physical node---> node num
    vnode_map_node    map[uint32] string  //vnode key ---->physical node
    
}


func ConHashNew() *ConHash {
    return &ConHash {
        hashfunc: crc32.ChecksumIEEE,
        hashring: HashRing{}, 
        vnodes: 0,
        node_map_num: make(map[string] int),
        vnode_map_node: make(map[uint32] string),
    }
}

func (ch *ConHash) NodeHash(data []byte) uint32 {
    return ch.hashfunc(data) 
}

//vnode format:  "n_name_%d"
func (ch *ConHash) NodeAdd(n_name string, vn_num int) {
    var name string
    var key uint32
    if _,ok := ch.node_map_num[n_name]; ok {
        fmt.Println("node", n_name, "already exist")
        return
    }
    ch.node_map_num[n_name] = vn_num
    for i := 0; i < vn_num; i++ {
        name = n_name + "_" + strconv.Itoa(i)
        key = ch.NodeHash([]byte(name))
        ch.vnode_map_node[key] = n_name
        ch.hashring = append(ch.hashring, key)
        ch.vnodes++
    }
    sort.Sort(ch.hashring)
}

func (ch *ConHash) NodeRemove(n_name string) {
    var name string
    var key uint32

    if _,ok := ch.node_map_num[n_name]; !ok {
       fmt.Println("node", n_name, "not exist") 
    }

    vnode_num := ch.node_map_num[n_name]
    for i := 0; i < vnode_num; i++ {
        name = n_name + "_" + strconv.Itoa(i)
        key = ch.NodeHash([]byte(name))
        for i,v := range ch.hashring {
            if v == key {
                ch.hashring = append(ch.hashring[:i], ch.hashring[i+1:]...)
            }
        }
        delete(ch.vnode_map_node, key)
        ch.vnodes--
    }
    delete(ch.node_map_num, n_name)
}


func (ch *ConHash) NodeLookup(object string) string {
    var key uint32
    var hitindex int

    key = ch.NodeHash([]byte(object))
    index := sort.Search(len(ch.hashring), func(i int) bool { return ch.hashring[i] >= key })
    if index == len(ch.hashring) {
        hitindex = 0
    } else {
        hitindex = index
    }
    hitkey := ch.hashring[hitindex]
    node, _ := ch.vnode_map_node[hitkey]
    return node
}

func (ch *ConHash) NodeGetVnodes() int {
    return ch.vnodes
}


func main() {
    var name string
    var nodestring string
    conhash := ConHashNew()
    fmt.Println("Add Node....")
    conhash.NodeAdd("titanic", 32)
    conhash.NodeAdd("terminator2018", 24)
    conhash.NodeAdd("Xenomorph", 25)
    conhash.NodeAdd("True Lies", 10)
    conhash.NodeAdd("avantar", 48)
    vnodes := conhash.NodeGetVnodes()
    fmt.Println("virtual nodes number: ",vnodes)
    fmt.Println("Lookup Node...") 
    for i := 0; i < 20; i++ {
        name = "shanghai_" + strconv.Itoa(i)
        nodestring = conhash.NodeLookup(name)
        fmt.Println(name, "match", nodestring)
    }
    fmt.Println("Delete Node...")
    conhash.NodeRemove("terminator2018")
    vnodes = conhash.NodeGetVnodes()
    fmt.Println("virtual nodes number: ",vnodes)
    for i := 0; i < 20; i++ {
        name = "shanghai_" + strconv.Itoa(i)
        nodestring = conhash.NodeLookup(name)
        fmt.Println(name, "match", nodestring)
    }
}
