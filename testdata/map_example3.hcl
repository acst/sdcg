package "github.com/acst/sdcg/testdata/package1" {}

use "github.com/acst/sdcg/testdata/package2" {}

map {
    from "privateStructA" {}
    to "package2.PublicStructA" {}
}

map {
    from "package2.PublicStructA" {}
    to "privateStructA" {}
}