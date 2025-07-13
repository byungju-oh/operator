# Go Operator (kubebuilder)

kubebuilder를 사용한 Kubernetes 오퍼레이터입니다.

## 🐹 특징

- **타입 안정성**: 컴파일 타임 에러 검출
- **고성능**: Go의 효율적인 리소스 사용
- **Kubernetes 네이티브**: controller-runtime 기반
- **교육용 구현**: 로그 출력으로 오퍼레이터 동작 학습

## 🚀 로컬 개발

### 사전 요구사항
```bash
# Go 1.19+ 설치
go version

# kubebuilder 설치 (선택사항)
curl -L -o kubebuilder https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

로컬 실행
bash# CRD 설치
make install

# 컨트롤러 로컬 실행
make run
테스트
bash# 다른 터미널에서
kubectl apply -f k8s/example-message.yaml
kubectl get messages
🔧 클러스터 배포
bash# Docker 이미지 빌드 및 배포
make docker-build docker-push IMG=gcr.io/YOUR_PROJECT_ID/go-operator:latest
make deploy IMG=gcr.io/YOUR_PROJECT_ID/go-operator:latest
📋 주요 기능

✅ Message CR 감지 및 로그 출력
✅ Reconcile 로직 구현
✅ 상태 관리
✅ 이벤트 기록