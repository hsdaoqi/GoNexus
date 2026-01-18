import grpc
from concurrent import futures
import time
import sys
import os

# 路径修复
sys.path.append(os.path.join(os.path.dirname(__file__), 'proto'))

import ai_service_pb2
import ai_service_pb2_grpc
from core.rag import RAGEngine

# 初始化引擎
engine = RAGEngine()


class AIService(ai_service_pb2_grpc.AIServiceServicer):

    def SyncMessage(self, request, context):
        try:
            user_id = f"User_{request.user_id}"
            # 获取昵称，没有则用ID
            nick = request.nickname if request.nickname else user_id

            engine.add_memory(request.content, user_id, request.msg_id, nick, request.session_id)
            return ai_service_pb2.SyncResponse(code=200)
        except Exception as e:
            print(f"❌ Error: {e}")
            return ai_service_pb2.SyncResponse(code=500)

    def RevokeMessage(self, request, context):
        try:
            engine.revoke_memory(request.msg_id)
            return ai_service_pb2.RevokeResponse(code=200)
        except Exception as e:
            print(f"❌ Error Revoke: {e}")
            return ai_service_pb2.RevokeResponse(code=500)

    def SemanticSearch(self, request, context):
        try:
            answer = engine.rag_qa(request.query, request.session_id, request.limit)
            return ai_service_pb2.SearchResponse(code=200, answer=answer)
        except Exception as e:
            return ai_service_pb2.SearchResponse(code=500, answer="服务异常")

    def ChatSummary(self, request, context):
        try:
            summary = engine.chat_summary(request.chats)
            return ai_service_pb2.SummaryResponse(code=200, summary=summary)
        except Exception as e:
            return ai_service_pb2.SummaryResponse(code=500)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    ai_service_pb2_grpc.add_AIServiceServicer_to_server(AIService(), server)
    server.add_insecure_port('[::]:50051')
    print("🚀 GoNexus-AI 服务已启动 (Port: 50051)...", flush=True)
    server.start()
    try:
        while True: time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
