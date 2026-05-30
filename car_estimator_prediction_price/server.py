import grpc
from concurrent import futures
import logging
from logging.handlers import RotatingFileHandler

import prediction_pb2
import prediction_pb2_grpc
from predictor.analyze import get_car_info, get_photos, get_sells, graph_build
logger = logging.getLogger('PredictionServer')
logger.setLevel(logging.INFO)

console_h = logging.StreamHandler()
console_h.setFormatter(logging.Formatter(
    '%(asctime)s %(levelname)8s %(name)s: %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
))
logger.addHandler(console_h)

file_h = RotatingFileHandler(
    filename='server.log',
    maxBytes=1_000_000,
    backupCount=5,
    encoding='utf-8'
)
file_h.setFormatter(logging.Formatter(
    '%(asctime)s %(levelname)8s %(name)s: %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
))
logger.addHandler(file_h)


class PredictionService(prediction_pb2_grpc.PredictionServiceServicer):
    def Predict(self, request, context):
        logger.info(
            "Received Predict request: make=%s, model=%s, year=%d, hp=%d, "
            "body=%s, yearsell=%d, odometer=%d, color=%s",
            request.make, request.model, request.year, request.hp,
            request.body, request.yearsell, request.odometer, request.color
        )

        try:
            price = get_car_info(
                request.make,
                request.model,
                request.year,
                request.hp,
                request.body,
                request.yearsell,
                request.odometer,
                request.color
            )
            photos = get_photos(
                request.make,
                request.model,
                request.year
            )
            sells = get_sells(
                request.make,
                request.model
            )
            graph_png = graph_build(
                request.make,
                request.model,
                request.year,
                request.hp,
                request.body,
                request.yearsell,
                request.odometer,
                request.color
            )

            response = prediction_pb2.PredictResponse(
                price=price,
                photo_urls=photos,
                sell_count=sells,
                graph_png=graph_png
            )
            logger.info(
                "Sending Predict response: price=%d, photos=%d urls, sells=%d",
                price, len(photos), sells
            )
            return response

        except Exception:
            logger.exception("Error handling Predict request")
            context.set_details('Internal server error')
            context.set_code(grpc.StatusCode.INTERNAL)
            return prediction_pb2.PredictResponse()
    
    def GetImages(self, request, context):
        try:
            logger.info(
                "Received request for images: make=%s, model=%s, year=%d",
                request.make, request.model, request.year
            )

            photos = get_photos(
                request.make,
                request.model,
                request.year
            )

            response = prediction_pb2.ImagesResponse(
                photo_urls=photos,
            )

            logger.info(
                "Sending images: photos=%d urls", len(photos)
            )

            return response
        
        except Exception:
            logger.exception("Error handling images request")
            context.set_details('Internal server error')
            context.set_code(grpc.StatusCode.INTERNAL)
            return prediction_pb2.ImagesResponse()



def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    prediction_pb2_grpc.add_PredictionServiceServicer_to_server(
        PredictionService(), server
    )

    port = '[::]:50051'
    server.add_insecure_port(port)
    server.start()
    logger.info("gRPC server started, listening on %s", port)
    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Server stopping by KeyboardInterrupt")
        server.stop(0)


if __name__ == '__main__':
    serve()