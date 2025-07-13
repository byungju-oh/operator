# Kubernetes Operators 학습 프로젝트

두 가지 언어와 프레임워크로 구현한 Kubernetes 오퍼레이터 학습 프로젝트입니다. 동일한 기능을 Python(kopf)과 Go(kubebuilder)로 각각 구현하여 두 접근 방식을 비교 학습할 수 있습니다.

## 📋 프로젝트 개요

이 프로젝트는 `Message` 커스텀 리소스를 관리하는 Kubernetes 오퍼레이터를 두 가지 방식으로 구현합니다:

- **Python 오퍼레이터**: kopf 프레임워크 사용, 실제 Pod 생성 및 관리
- **Go 오퍼레이터**: kubebuilder 사용, 로그 출력 및 상태 관리

## 🏗️ 프로젝트 구조

```
├── python-operator/          # Python + kopf 구현
│   ├── main.py              # 오퍼레이터 메인 로직
│   ├── Dockerfile           # 컨테이너 이미지
│   ├── requirements.txt     # Python 의존성
│   └── k8s/                 # Kubernetes 매니페스트
│       ├── crd.yaml
│       ├── operator-deploy.yaml
│       └── example-message.yaml
│
└── go-operator/             # Go + kubebuilder 구현
    ├── main.go              # 애플리케이션 진입점
    ├── Dockerfile           # 컨테이너 이미지
    ├── Makefile             # 빌드 및 배포 스크립트
    ├── api/v1/              # CRD 타입 정의
    ├── controllers/         # 컨트롤러 로직
    └── k8s/                 # Kubernetes 매니페스트
        ├── crd.yaml
        ├── operator-deploy.yaml
        └── example-message.yaml
```

## 🎯 핵심 기능

### Message CRD (Custom Resource Definition)
두 오퍼레이터 모두 동일한 `Message` CRD를 사용합니다:

```yaml
apiVersion: myorg.dev/v1
kind: Message
metadata:
  name: example-message
spec:
  text: "안녕하세요! 오퍼레이터 테스트 메시지입니다!"
```

### 구현별 차이점

| 기능 | Python 오퍼레이터 | Go 오퍼레이터 |
|------|------------------|---------------|
| **프레임워크** | kopf | kubebuilder |
| **주요 동작** | Pod 생성/관리 | 로그 출력/상태 관리 |
| **이벤트 처리** | 데코레이터 기반 | Reconcile 루프 |
| **개발 속도** | 빠름 (Python) | 보통 (타입 안정성) |
| **성능** | 보통 | 높음 (Go) |
| **학습 목적** | 실제 리소스 관리 | 오퍼레이터 패턴 이해 |

## 🚀 빠른 시작

### 사전 요구사항
- Kubernetes 클러스터 (minikube, kind, GKE 등)
- kubectl 설정 완료
- Docker (컨테이너 이미지 빌드용)

### 1. Python 오퍼레이터 실행

```bash
cd python-operator

# 로컬 개발 (권장)
pip install -r requirements.txt
kubectl apply -f k8s/crd.yaml
kopf run main.py --standalone

# 또는 클러스터 배포
kubectl apply -f k8s/
```

### 2. Go 오퍼레이터 실행

```bash
cd go-operator

# 로컬 개발 (권장)
make install  # CRD 설치
make run      # 로컬 실행

# 또는 클러스터 배포
make deploy IMG=gcr.io/YOUR_PROJECT_ID/go-operator:latest
```

### 3. 테스트

```bash
# Message 리소스 생성
kubectl apply -f python-operator/k8s/example-message.yaml

# 결과 확인
kubectl get messages
kubectl get pods  # Python 오퍼레이터의 경우
kubectl logs -f deployment/message-operator  # 오퍼레이터 로그
```

## 📚 학습 가이드

### 1. Python 오퍼레이터 학습 포인트
- **간단한 시작**: 데코레이터 기반의 직관적인 이벤트 핸들링
- **실제 동작**: Pod 생성/삭제를 통한 실질적인 리소스 관리 학습
- **빠른 프로토타이핑**: Python의 간결함으로 빠른 개발 사이클

```python
@kopf.on.create('myorg.dev', 'v1', 'messages')
def create_fn(spec, name, namespace, logger, **kwargs):
    text = spec.get('text', 'Hello from Operator!')
    create_pod(name, namespace, text, logger)
```

### 2. Go 오퍼레이터 학습 포인트
- **Kubernetes 네이티브**: controller-runtime을 사용한 표준적인 접근
- **타입 안정성**: 컴파일 타임 에러 검출로 안정적인 개발
- **프로덕션 준비**: 메트릭, 헬스체크, 리더 일렉션 등 내장

```go
func (r *MessageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Reconciliation 로직
    logger.Info("🔄 Reconciling Message", "name", message.Name)
    return ctrl.Result{}, nil
}
```

## 🔧 고급 사용법

### 개발 워크플로우

1. **로컬 개발**: 두 오퍼레이터 모두 클러스터 외부에서 로컬 실행 가능
2. **디버깅**: 상세한 로그와 이벤트를 통한 디버깅
3. **테스트**: Message 리소스 생성/수정/삭제로 동작 확인
4. **배포**: 컨테이너 이미지 빌드 후 클러스터 배포

### 커스터마이제이션

- **메시지 처리 로직 변경**: `main.py` 또는 `message_controller.go` 수정
- **새로운 필드 추가**: CRD 스키마와 타입 정의 확장
- **추가 리소스 관리**: Service, ConfigMap 등 다른 리소스 관리 로직 추가

## 🛡️ 보안 고려사항

두 오퍼레이터 모두 다음 보안 모범 사례를 적용했습니다:

- **최소 권한 원칙**: 필요한 최소한의 RBAC 권한만 부여
- **비루트 실행**: 컨테이너를 비루트 사용자로 실행
- **보안 컨텍스트**: 권한 상승 방지, seccomp 프로필 적용
- **리소스 제한**: CPU/메모리 리소스 제한 설정

## 🐛 문제 해결

### 일반적인 문제

1. **CRD 등록 실패**
   ```bash
   kubectl get crd messages.myorg.dev
   kubectl apply -f k8s/crd.yaml
   ```

2. **권한 오류**
   ```bash
   kubectl get clusterrolebinding | grep operator
   kubectl apply -f k8s/operator-deploy.yaml
   ```

3. **로그 확인**
   ```bash
   kubectl logs -f deployment/message-operator
   kopf run main.py --verbose  # Python 오퍼레이터 디버그 모드
   ```

## 🤝 기여하기

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📖 추가 학습 자료

- [Kubernetes Operators 패턴](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [kopf 공식 문서](https://kopf.readthedocs.io/)
- [kubebuilder 공식 문서](https://book.kubebuilder.io/)
- [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다. 자세한 내용은 `LICENSE` 파일을 참조하세요.

---

**🎯 학습 목표**: 이 프로젝트를 통해 Kubernetes 오퍼레이터의 핵심 개념을 이해하고, Python과 Go 두 가지 접근 방식의 장단점을 직접 체험해보세요!
