package "github.com/acst/sdcg/testdata/package3" {
    filename = "mapper"
}

use "github.com/acst/sdcg/testdata/package1" {
    alias = "something"
}

use "github.com/acst/sdcg/testdata/package2" {}

map {
    from "something.StructA" {
        ignore = [
            "Field1"
        ]
    }
    to "package2.StructC" {
        ignore = [
            "Field1"
        ]
    }
    field "Field2" {
        to = "Field4"
    }
    field "Field4" {
        to = "Field2"
    }
    field "Field9" {
        using = "${converter(sliceToString)}"
    }
}

map {
    from "something.StructB" {}
    to  "package2.StructD" {}
    field "Field1" {
        to = "Field2"
    }
    field "Field2" {
        to = "Field1"
    }
}

map {
    from "package2.StructC" {}
    to "something.StructA" {}
    field "Field2" {
        to = "Field4"
    }
    field "Field4" {
        to = "Field2"
    }
    field "Field9" {
        using = "${converter(stringToSlice)}"
    }
}

map {
    from "package2.StructD" {}
    to "something.StructB" {}
    field "Field1" {
        to = "Field2"
    }
    field "Field2" {
        to = "Field1"
    }
}

convert {
    type = "something.Enum"
    to = "${gotype(string)}"
    using = "enumToString"
}

convert {
    type = "${gotype(string)}"
    to = "something.Enum"
    using = "stringToEnum"
}