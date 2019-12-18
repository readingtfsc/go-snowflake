## snowflake for go

快速开始
安装
    go get github.com/night-reading/go-snowflake
例子
    func main() {
	    node := snowflake.NewNode(1)

	    id := node.Snowflake()
	    fmt.Println(id)
    }