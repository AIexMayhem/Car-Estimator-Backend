import grpc
import prediction_pb2
import prediction_pb2_grpc

def run():
    channel = grpc.insecure_channel('localhost:50051')
    stub = prediction_pb2_grpc.PredictionServiceStub(channel)

    request = prediction_pb2.PredictRequest(
        make="Audi",
        model="A8",
        year=2010,
        hp=300,
        body="Sedan",
        yearsell=2021,
        odometer=100000,
        color="White"
    )

    try:
        response = stub.Predict(request, timeout=5)
    except grpc.RpcError as e:
        print(f"gRPC error: {e.code()} – {e.details()}")
        return

    print(f"Price: {response.price}")
    print(f"Sell count: {response.sell_count}")
    print("Photo URLs:")
    for url in response.photo_urls:
        print("  ", url)
    with open("forecast.png", "wb") as f:
        f.write(response.graph_png)
        print("Graph written to forecast.png")

if __name__ == "__main__":
    run()