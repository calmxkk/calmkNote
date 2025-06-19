[toc]

# cloud-init

## 启动流程分析

cloud-init的执行包括5个阶段，执行阶段从前到后分别为：Generator、Local、Network、Config、Final，

### generator

- systemd服务：/usr/lib/systemd/system-generators/cloud-init-generator