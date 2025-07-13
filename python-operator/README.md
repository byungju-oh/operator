---

### `/python-operator/README.md`
```markdown
# Python Operator (kopf)

kopf 프레임워크를 사용한 Kubernetes 오퍼레이터입니다.

##  특징

- **빠른 개발**: Python의 간결한 문법으로 빠른 프로토타이핑
- **강력한 기능**: 데코레이터 기반의 직관적인 이벤트 핸들링
- **풍부한 생태계**: Python 라이브러리 활용 가능
- **실제 Pod 생성**: Message CR 생성 시 실제 Pod를 생성하고 관리

## 🚀 로컬 개발

### 사전 요구사항
```bash
pip install -r requirements.txt

로컬 실행
bash# CRD 등록 (한 번만)
kubectl apply -f k8s/crd.yaml

# 오퍼레이터 로컬 실행
kopf run main.py --standalone
테스트
bash# 다른 터미널에서
kubectl apply -f k8s/example-message.yaml
kubectl get pods
kubectl logs msg-hello-msg
🔧 클러스터 배포
bash# 전체 배포
kubectl apply -f k8s/

# 또는 개별 배포
kubectl apply -f k8s/crd.yaml
kubectl apply -f k8s/operator-deploy.yaml
kubectl apply -f k8s/example-message.yaml
📋 주요 기능

✅ Message CR 생성 시 자동 Pod 생성
✅ Message 내용 변경 시 자동 업데이트
✅ Message 삭제 시 자동 정리
✅ 실시간 로그 출력
✅ 에러 핸들링