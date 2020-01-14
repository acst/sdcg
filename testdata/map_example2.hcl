package "github.com/acst/sdcg/testdata/package1"{}

map {
    from "privateStructA" {}
    to "privateStructB" {}
}

map {
    from "privateStructB" {}
    to "privateStructA" {}
}