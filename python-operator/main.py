import kopf
import kubernetes.client
from kubernetes.client import V1Pod, V1PodSpec, V1Container, V1ObjectMeta
import asyncio

@kopf.on.create('myorg.dev', 'v1', 'messages')
def create_fn(spec, name, namespace, logger, **kwargs):
    text = spec.get('text', 'Hello from Operator!')
    logger.info(f"🆕 [CREATE] Message '{name}' 생성됨: {text}")
    
    create_pod(name, namespace, text, logger)

@kopf.on.update('myorg.dev', 'v1', 'messages')  
def update_fn(spec, name, namespace, old, new, logger, **kwargs):
    old_text = old.get('spec', {}).get('text', '')
    new_text = spec.get('text', '')
    
    if old_text != new_text:
        logger.info(f"🔄 [UPDATE] Message '{name}' 변경: '{old_text}' → '{new_text}'")
        delete_pod(name, namespace, logger)
        create_pod(name, namespace, new_text, logger)

@kopf.on.delete('myorg.dev', 'v1', 'messages')
def delete_fn(name, namespace, logger, **kwargs):
    logger.info(f"🗑️ [DELETE] Message '{name}' 삭제됨")
    delete_pod(name, namespace, logger)

def create_pod(name, namespace, text, logger):
    escaped_text = text.replace('"', '\\"').replace("'", "\\'")
    
    container = V1Container(
        name="echo",
        image="busybox",
        command=["/bin/sh", "-c", f'while true; do echo "[$(date)] {escaped_text}"; sleep 10; done'],
        resources=kubernetes.client.V1ResourceRequirements(
            requests={"memory": "64Mi", "cpu": "50m"},
            limits={"memory": "128Mi", "cpu": "100m"}
        ),
        security_context=kubernetes.client.V1SecurityContext(
            allow_privilege_escalation=False,
            run_as_non_root=True,
            run_as_user=1000,
            capabilities=kubernetes.client.V1Capabilities(drop=["ALL"]),
            seccomp_profile=kubernetes.client.V1SeccompProfile(type="RuntimeDefault")
        )
    )
    
    pod = V1Pod(
        metadata=V1ObjectMeta(
            name=f"msg-{name}",
            labels={"app": f"msg-{name}", "managed-by": "message-operator"}
        ),
        spec=V1PodSpec(
            containers=[container], 
            restart_policy="Never",
            security_context=kubernetes.client.V1PodSecurityContext(
                run_as_non_root=True,
                run_as_user=1000
            )
        )
    )
    
    try:
        api = kubernetes.client.CoreV1Api()
        api.create_namespaced_pod(namespace=namespace, body=pod)
        logger.info(f"✅ Pod 'msg-{name}' 생성됨")
    except Exception as e:
        logger.error(f"❌ Pod 생성 실패: {e}")

def delete_pod(name, namespace, logger):
    try:
        api = kubernetes.client.CoreV1Api()
        api.delete_namespaced_pod(name=f"msg-{name}", namespace=namespace)
        logger.info(f"✅ Pod 'msg-{name}' 삭제됨")
    except Exception as e:
        logger.warning(f"⚠️ Pod 삭제 실패: {e}")