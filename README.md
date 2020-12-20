# GraphQL 和 GORM

本文档结合了 GraphQL 和 GORM 库，目的实现基于 MySQL 数据库的 WEB API 框架。

首先新建一个空目录,存放我们的项目

```shell
mkdir go-orders-graph-api/graph -pv
cd go-orders-graph-api
```



## 定义数据结构

首先定义一个数据库所用的数据结构。创建一个文件**`graph/schema.graphqls`**

```
type Order {
    id: Int!
    customerName: String!
    orderAmount: Float!
    items: [Item!]!
}

type Item {
    id: Int!
    productCode: String!
    productName: String!
    quantity: Int!
}

input OrderInput {
    customerName: String!
    orderAmount: Float!
    items: [ItemInput!]!
}

input ItemInput {
    productCode: String!
    productName: String!
    quantity: Int!
}

type Mutation {
    createOrder(input: OrderInput!): Order!
    updateOrder(orderId: Int!, input: OrderInput!): Order!
    deleteOrder(orderId: Int!): Boolean!
}

type Query {
    orders: [Order!]!
}
```



## gqlgen 工具

使用 `gqlgen` 工具可以依据你的数据结构也就是上面的文件`graph/schema.graphqls`生成几乎除了逻辑外的大多数代码，设置包括了使用`net/http`库启动的服务端。

1. 首先需要获得`gqlgen`以及数据库开发所需的`gorm`工具库

```shell
shell$ go get -u github.com/jinzhu/gorm
shell$ go get -u github.com/go-sql-driver/mysql
shell$ go get github.com/99designs/gqlgen
```

 2.  进入到项目的目录下,初始化go项目的模块

```shell
shell$ cd go-orders-graph-api
shell$ go mod init github.com/yangwawa0323/go-orders-graph-api
```

 3.  接着初始化形成

```shell
shell$ go run github.com/99designs/gqlgen init
```

或者直接使用编译后的命令

```shell
shell$ gqlgen init
```

4. 这样就形成下面的文件结构

```
├── go.mod
├── go.sum
├── gqlgen.yml               - gqlgen 配置文件,用来控制生成的代码.
├── graph
│   ├── generated            - 仅仅包含生成的generated包
│   │   └── generated.go
│   ├── model                - 生成的所有数据模型
│   │   └── models_gen.go
│   ├── resolver.go          - 此文件不会因为重新生成而覆盖
│   ├── schema.graphqls      - 数据结构
│   └── schema.resolvers.go  - 自己定义的如何实现 schema.graphql逻辑功能
└── server.go                - Web启动服务
  
```



## graph/schema.graphqls结构文件

1. 默认`qglgen`将会自动生成一个Todo的案例，我们需要删除以前的内容

```shell
shell$ rm ./graph/generated/* ./graph/model/* 
```

然后可以开始重新定义数据结构,清空之前的Todo结构，我们将建立一个订单的范例

```
shell$ vim graph/schema.graphqls

type Order {
    id: Int!
    customerName: String!
    orderAmount: Float!
    items: [Item!]!
}

type Item {
    id: Int!
    productCode: String!
    productName: String!
    quantity: Int!
}

input OrderInput {
    customerName: String!
    orderAmount: Float!
    items: [ItemInput!]!
}

input ItemInput {
    productCode: String!
    productName: String!
    quantity: Int!
}

type Mutation {
    createOrder(input: OrderInput!): Order!
    updateOrder(orderId: Int!, input: OrderInput!): Order!
    deleteOrder(orderId: Int!): Boolean!
}

type Query {
    orders: [Order!]!
}
```





## 模型的定义

一旦修改好，再次生成我们的代码

```shell
gqlgen generate
```

这样我们再次得到之前的目录结构中的每个文件。在整个交互中，我们只需保留`schema.graphqls`文件即可

生成的模型`graph/models/models_gen.go`如下

```go
package model
//...
type Item struct {
	ID          int `json:"id"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

type Order struct {
	ID           int  `json:"id"`
	CustomerName string  `json:"customerName"`
	OrderAmount  float64 `json:"orderAmount"`
	Items        []*Item `json:"items"`
}
```

为了实现GORM数据结构，我们可以添加关于外键的定义,按照需求,一个订单有多条记录，添加注解`gorm:foreignkey:ID`到`Order`下的`Items`,这将以为着`Item`表中的`ID`列外键形式参考的为`Order`表的`ID`，详细请[参考GORM](https://gorm.io/docs/has_many.html)

```go
package model
//...
type Order struct {
	ID           int  `json:"id"`
	CustomerName string  `json:"customerName"`
	OrderAmount  float64 `json:"orderAmount"`
    Items        []*Item `json:"items" gorm:"foreignkey:ID"`
}
```

最终，如下

```go
package model

type ItemInput struct {
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

type OrderInput struct {
	CustomerName string       `json:"customerName"`
	OrderAmount  float64      `json:"orderAmount"`
	Items        []*ItemInput `json:"items"`
}

type Item struct {
	ID          int `json:"id"`
	ProductCode string `json:"productCode"`
	ProductName string `json:"productName"`
	Quantity    int    `json:"quantity"`
}

type Order struct {
	ID           int  `json:"id"`
	CustomerName string  `json:"customerName"`
	OrderAmount  float64 `json:"orderAmount"`
    Items        []*Item `json:"items" gorm:"foreignkey:ID"`
}
```



## 导入开发库

思路：首先我们需要用到数据库，因此我们在以下文件中导入我们开发所用的 `gorm`库

```go
// In resolver.go
import "github.com/jinzhu/gorm"
```

> 注意：案例中代码使用到的非`grom.io/grom`,而是`github.com`下的。请保持一致，**混杂使用两种库，一定会出现看上去切片指针一样，而在使用中总是报错的现象**
>
> 还有一点，如果你使用`vscode`编写代码`golint`会在保存时**自动去除没有使用但导入的库**，你所做的操作就会白费。

```go
// In schema.resolver.go
import (
	"context"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/generated"
	"github.com/yangwawa0323/go-orders-graphql-api/graph/model"
)
```



```go
// In server.go
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/soberkoder/go-orders-graphql-api/graph"
	"github.com/soberkoder/go-orders-graphql-api/graph/generated"
	"github.com/soberkoder/go-orders-graphql-api/graph/model"
)
```

> 注意：引入模块的名称是基于`go mod init`建立项目时的名称。如果你是复制以上代码，请自行修正



## 初始化数据库

为了将 Graphql 和 GORM结合，我们可以在项目根目录下的服务启动程序`server.go`中添加初始化数据库代码

```go
# In server.go
var db *gorm.DB;

func initDB() {
    var err error
    dataSourceName := "root:redhat@tcp(localhost:3306)/?parseTime=True&charset=utf8mb4"
    db, err = gorm.Open("mysql", dataSourceName)

    if err != nil {
        fmt.Println(err)
        panic("failed to connect database")
    }

    db.LogMode(true)

    // Create the database. This is a one-time step.
    // Comment out if running multiple times - You may see an error otherwise
    db.Exec("CREATE DATABASE IF NOT EXISTS test_db")
    db.Exec("USE test_db")

    // Migration to create tables for Order and Item schema
    db.AutoMigrate(&models.Order{}, &models.Item{})	
}
```



## GORM 数据库连接

resolver.go 文件中添加初始化数据库db的指针

```go
# In resolver.go
type Resolver struct{
    DB *gorm.DB
}
```

在服务器中调用先前的  **initDB** 函数， 并在 Resolver 初始化的传递已经建立连接的 DB 对象

```go
# In server.go
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	initDB()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
```



## 双向数据操作

之前在模型定义中的`createOrder`创建订单、`updateOrder`更新订单、`deleteOrder`删除订单的逻辑都没有实现，因此这是整个程序编写中除了模型定义，初始化数据库连接外，我们需要自主实现的代码编写部分

1. 创建订单逻辑，新建一个Order对象，使用DB.Create保存到数据库，再将订单数据返回，交由`generated`库序列化转换成`JSON`格式返回给客户端

```
func (r *mutationResolver) CreateOrder(ctx context.Context, input OrderInput) (*models.Order, error) {
    order := models.Order {
        CustomerName: input.CustomerName,
        OrderAmount: input.OrderAmount,
        Items: mapItemsFromInput(input.Items),
    }
    r.DB.Create(&order)
    return &order, nil
}
```

2. 更新订单，按照订单ID号，以及数据建立新的对象，使用DB.Save保存到数据库，同样返回给`generated`库序列化成`JSON`结构返回给客户

```go
func (r *mutationResolver) UpdateOrder(ctx context.Context, orderID int, input OrderInput) (*models.Order, error) {
    updatedOrder := models.Order {
        ID: orderID,
        CustomerName: input.CustomerName,
        OrderAmount: input.OrderAmount,
        Items: mapItemsFromInput(input.Items),
    }
    r.DB.Save(&updatedOrder)
    return &updatedOrder, nil
}
```

3. 删除订单，从数据库中按订单ID找出订单让后使用 DB.Delete 删除从表和主表的记录。

```go
func (r *mutationResolver) DeleteOrder(ctx context.Context, orderID int) (bool, error) {
    r.DB.Where("id = ?", orderID).Delete(&models.Item{})
    r.DB.Where("id = ?", orderID).Delete(&models.Order{})
    return true, nil;
}
```



## 查询操作

```go
func (r *queryResolver) Orders(ctx context.Context) ([]*models.Order, error) {	
    var orders []*models.Order
    r.DB.Preload("Items").Find(&orders)
    
    return orders, nil
}
```



## 运行服务

```go
go run server.go
```

打开浏览器访问 **http://localhost:8080**



## 测试

* **创建订单**

```
mutation createOrder ($input: OrderInput!) {
  createOrder(input: $input) {
    id
    customerName
    items {
      id
      productCode
      productName
      quantity
    }
  }
}
```

传入的参数

```json
{
  "input": {
    "customerName": "Leo",
    "orderAmount": 9.99,
    "items": [
      {
      "productCode": "2323",
      "productName": "IPhone X",
      "quantity": 1
      }
    ]
  }
}
```

* 查询订单

```
  query orders {
    orders {
      id  
      customerName
      items {
        productName
        quantity
      }
    }
  }
```

  

* 更新订单

```
  mutation updateOrder ($orderId: Int!, $input: OrderInput!) {
    updateOrder(orderId: $orderId, input: $input) {
      id
      customerName
      items {
        id
        productCode
        productName
        quantity
      }
    }
  }
```

  传入参数

```json
  {
    "orderId":1,
    "input": {
      "customerName": "Cristiano",
      "orderAmount": 9.99,
      "items": [
        {
        "productCode": "2323",
        "productName": "IPhone X",
        "quantity": 1
        }
      ]
    }
  }
```

  

* 删除订单

```
mutation deleteOrder ($orderId: Int!) {
    deleteOrder(orderId: $orderId)
}
```

  传入参数

```json
{
  "orderId": 3
}
```

  
