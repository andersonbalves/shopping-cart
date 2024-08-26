# diagram.py
from diagrams import Diagram
from diagrams.aws.compute import Lambda, LambdaFunction
from diagrams.aws.network import APIGateway, APIGatewayEndpoint
from diagrams.aws.database import Dynamodb
from diagrams.aws.management import Cloudwatch
from diagrams.onprem.database import MySQL
from diagrams import Cluster, Diagram, Edge
from diagrams.programming.framework import Angular


with Diagram("Shopping Cart", show=False, filename="shopping-cart", direction="LR"):
    with Cluster("Docker"):
        with Cluster("Gateway"):
            products_api_gateway = APIGateway("products")
            cart_api_gateway = APIGateway("shopping-cart")

        with Cluster("Functions"):
            products_lambda = Lambda("ProductsAPI")
            cart_lambda = Lambda("ShoppingCartAPI")

        with Cluster("Database"):
            products_sql_base = MySQL("Products")
            cart_dynamodb = Dynamodb("ShoppingCart")

        cloud_watch = Cloudwatch("Logs")

    front = Angular("Shopping")

    front >> Edge(label="GET") >> products_api_gateway >> products_lambda

    products_lambda >> Edge(label="SELECT ALL") >> products_sql_base
    products_lambda >> Edge(style="dashed") >> cloud_watch

    front >> Edge(label="GET/POST/PUT/DELETE") >> cart_api_gateway >> cart_lambda
    cart_lambda >> Edge(label="CRUD") >> cart_dynamodb
    cart_lambda >> Edge(style="dashed") >> cloud_watch
