# Kubernetes 오퍼레이터란 무엇인가? 🤖

이 프로젝트는 **GKE(Google Kubernetes Engine)** 환경에서 Kubernetes 오퍼레이터의 핵심 개념을 학습하기 위한 실습용 코드입니다. Python과 Go 두 가지 방식으로 동일한 기능을 구현하여 오퍼레이터 패턴을 깊이 이해할 수 있습니다.
이 프로젝트는 `Message` 커스텀 리소스를 관리하는 Kubernetes 오퍼레이터를 두 가지 방식으로 구현합니다:

- **Python 오퍼레이터**: kopf 프레임워크 사용, 실제 Pod 생성 및 관리
- **Go 오퍼레이터**: kubebuilder 사용, 로그 출력 및 상태 관리

> **"오퍼레이터는 Kubernetes에서 애플리케이션을 자동으로 관리하는 소프트웨어입니다"**


### 1. 오퍼레이터의 정의
**오퍼레이터는 인간 운영자(Human Operator)의 지식을 코드로 구현한 것입니다.**

```
전통적 방식: 인간이 수동으로 관리
kubectl create deployment myapp
kubectl scale deployment myapp --replicas=5
kubectl delete pod myapp-xxx  # 문제 발생 시

오퍼레이터 방식: 자동화된 관리
apiVersion: myorg.dev/v1
kind: Message  # 원하는 상태만 선언
```

### 2. 오퍼레이터의 핵심 개념

#### 📝 **CRD (Custom Resource Definition)**
Kubernetes에 새로운 리소스 타입을 정의합니다.
```yaml
# 우리가 만든 새로운 리소스 타입
apiVersion: myorg.dev/v1
kind: Message
spec:
  text: "안녕하세요!"
```

#### 🔄 **컨트롤 루프 (Control Loop)**
```
현재 상태 → 원하는 상태 비교 → 조치 → 다시 확인 → 반복
```

#### 🎛️ **컨트롤러 (Controller)**
실제로 상태를 관리하는 코드입니다.

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


## 🏗️ 이 프로젝트에서 배우는 것

### Python 오퍼레이터 (kopf) - 실제 리소스 관리 학습
```python
@kopf.on.create('myorg.dev', 'v1', 'messages')
def create_fn(spec, name, namespace, logger, **kwargs):
    # Message가 생성되면 → Pod를 자동 생성
    text = spec.get('text', 'Hello!')
    create_pod(name, namespace, text, logger)
```

**학습 포인트:**
- ✅ Message 생성 → Pod 자동 생성
- ✅ Message 수정 → Pod 자동 업데이트  
- ✅ Message 삭제 → Pod 자동 정리
- ✅ 실제 리소스가 어떻게 관리되는지 체험

### Go 오퍼레이터 (kubebuilder) - 오퍼레이터 패턴 학습
```go
func (r *MessageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Reconcile 패턴 구현
    logger.Info("🔄 Reconciling Message", "name", message.Name)
    
    // 상태 업데이트
    message.Status.Phase = "Processed"
    return ctrl.Result{}, nil
}
```

**학습 포인트:**
- ✅ Reconcile 패턴의 이해
- ✅ 상태(Status) 관리 방법
- ✅ 이벤트 기록과 모니터링
- ✅ Kubernetes 네이티브 패턴

## 🚀 GKE에서 실습하기

### 사전 준비
```bash
# GKE 클러스터가 준비되어 있다고 가정
gcloud container clusters get-credentials [클러스터명] --zone [존] --project [프로젝트ID]
kubectl get nodes  # 클러스터 연결 확인
```

### 1단계: 오퍼레이터 없이 수동 관리 (Before)
```bash
# 전통적인 방식 - 모든 것을 수동으로
kubectl create deployment nginx --image=nginx
kubectl scale deployment nginx --replicas=3
kubectl delete pod nginx-xxx  # 문제 발생 시 수동 대응
```

### 2단계: Python 오퍼레이터로 자동 관리 (After)
```bash
cd python-operator

# 1. CRD 설치 (새로운 리소스 타입 등록)
kubectl apply -f k8s/crd.yaml
kubectl get crd messages.myorg.dev  # 확인

# 2. 오퍼레이터 실행 (로컬)
pip install -r requirements.txt
kopf run main.py --standalone

# 3. 다른 터미널에서 테스트
kubectl apply -f k8s/example-message.yaml

# 4. 마법 같은 일이 벌어집니다! 🎪
kubectl get messages
kubectl get pods  # Pod가 자동으로 생성됨!
kubectl logs msg-hello-msg  # 메시지가 계속 출력됨
```

### 3단계: Go 오퍼레이터로 패턴 이해
```bash
cd go-operator

# 1. CRD 설치
make install

# 2. 컨트롤러 실행
make run

# 3. 다른 터미널에서 테스트
kubectl apply -f k8s/example-message.yaml
kubectl get messages -o wide  # 상태 확인
```

## 🔍 실습으로 배우는 오퍼레이터 동작 원리

### 실험 1: 생성 이벤트 관찰
```bash
# Message 생성
kubectl apply -f - <<EOF
apiVersion: myorg.dev/v1
kind: Message
metadata:
  name: test-msg
spec:
  text: "실습용 메시지입니다!"
EOF

# 결과 관찰
kubectl get pods | grep test-msg      # Python: Pod 생성됨
kubectl get messages test-msg -o yaml # Go: 상태 업데이트됨
```

### 실험 2: 수정 이벤트 관찰
```bash
# 메시지 내용 변경
kubectl patch message test-msg --type='merge' -p='{"spec":{"text":"변경된 메시지!"}}'

# Python 오퍼레이터: 기존 Pod 삭제 후 새 Pod 생성
# Go 오퍼레이터: Reconcile 이벤트 발생
```

### 실험 3: 삭제 이벤트 관찰
```bash
kubectl delete message test-msg

# Python: Pod도 함께 정리됨 (Garbage Collection)
# Go: 삭제 이벤트 로그 출력
```

## 💡 오퍼레이터 패턴의 핵심 이해

### 선언적 API (Declarative API)
```yaml
# "어떻게"가 아닌 "무엇을" 원하는지만 선언
spec:
  text: "이런 메시지를 원합니다"
# 오퍼레이터가 알아서 "어떻게" 처리함
```

### 이벤트 기반 아키텍처
```
Message 생성 → 이벤트 발생 → 오퍼레이터 반응 → Pod 생성
Message 수정 → 이벤트 발생 → 오퍼레이터 반응 → Pod 업데이트
Message 삭제 → 이벤트 발생 → 오퍼레이터 반응 → Pod 정리
```

### 지속적인 조정 (Continuous Reconciliation)
```
원하는 상태: Message{text: "Hello"}
현재 상태: Pod가 없음
조치: Pod 생성
결과: 상태 일치 ✅

원하는 상태: Message{text: "Hi"}  
현재 상태: Pod{text: "Hello"}
조치: Pod 재생성
결과: 상태 일치 ✅
```


## 🤔 왜 오퍼레이터를 사용하나요?

### Before: 수동 관리의 문제점
- ❌ 24/7 인력 필요
- ❌ 실수 가능성
- ❌ 일관성 부족
- ❌ 확장성 한계

### After: 오퍼레이터의 장점
- ✅ 자동화된 관리
- ✅ 일관된 동작
- ✅ 도메인 지식 코드화
- ✅ 확장 가능한 관리



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
