ARG GO_VER=1.22
ARG GO_ARCH=alpine

ARG INSTALL_KUBECTL=true
ARG KUBECTL_VER=1.29.1

ARG INSTALL_KUBEBUILDER=false
ARG KUBEBUILDER_VER=3.10.0

FROM golang:${GO_VER}-${GO_ARCH}

# https://blog.scottchayaa.com/post/2018/11/04/docker-arg-env-variable/#arg
# 1. ARG need to be declared again after FROM
# 2. To get benfit of ARG, you need to put ARG as close as to where it is used.
#    If you change value when building, it will invalidate cache and rebuild.


# Install Golang extension package tools
RUN go install github.com/ramya-rao-a/go-outline@latest         && \
    go install github.com/cweill/gotests/gotests@latest         && \
    go install github.com/fatih/gomodifytags@latest             && \
    go install github.com/josharian/impl@latest                 && \
    go install github.com/haya14busa/goplay/cmd/goplay@latest   && \
    go install github.com/go-delve/delve/cmd/dlv@latest         && \
    go install honnef.co/go/tools/cmd/staticcheck@latest        && \
    go install golang.org/x/tools/gopls@latest

ARG GO_ARCH
RUN if [ "${GO_ARCH}" == "alpine" ]; then \
        apk add git curl bash alpine-sdk bash-completion ncurses yq helm && \
        sed -i 's/\/bin\/sh/\/bin\/bash/g' /etc/passwd; \
    fi

# Install basic tools
RUN echo "alias ll='ls -al'"                         >> /root/.bashrc

ARG KUBECTL_VER
ARG INSTALL_KUBECTL
# Install kubectl
RUN if [ "${INSTALL_KUBECTL}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}')                                  && \
        curl -LO "https://dl.k8s.io/release/v${KUBECTL_VER}/bin/${OS}/${ARCH}/kubectl"  && \
        chmod +x kubectl && mv kubectl /usr/local/bin/kubectl                           && \
        echo 'source <(kubectl completion bash)'                    >> /root/.bashrc    && \
        echo 'alias k=kubectl'                                      >> /root/.bashrc    && \
        echo 'complete -o default -F __start_kubectl k'             >> /root/.bashrc    && \
        git -C /opt clone https://github.com/jonmosco/kube-ps1.git                      && \
        echo 'source /opt/kube-ps1/kube-ps1.sh'                     >> /root/.bashrc    && \
        echo "PS1='[\u@\h \W \$(kube_ps1)]\$ '"                     >> /root/.bashrc    && \
        echo "KUBE_PS1_SYMBOL_COLOR=''"                             >> /root/.bashrc    && \
        echo "KUBE_PS1_CTX_COLOR=''"                                >> /root/.bashrc    && \
        echo "KUBE_PS1_NS_COLOR=''"                                 >> /root/.bashrc    && \
        echo "KUBE_PS1_SYMBOL_ENABLE='false'"                       >> /root/.bashrc    && \
        echo "KUBE_PS1_ENABLED='off'"                               >> /root/.bashrc    && \
        cd "$(mktemp -d)"                                           && \
        KREW="krew-${OS}_${ARCH}"                                   && \
        curl -fsSLO "https://github.com/kubernetes-sigs/krew/releases/latest/download/${KREW}.tar.gz" && \
        tar zxvf "${KREW}.tar.gz"                                   && \
        echo 'export PATH="${KREW_ROOT:-$HOME/.krew}/bin:$PATH"'    >> /root/.bashrc    && \
        ./"${KREW}" install krew                                    && \
        ./"${KREW}" install ctx                                     && \
        ./"${KREW}" install ns                                      && \
        ./"${KREW}" install neat                                    ;\
    fi


ARG INSTALL_KUBEBUILDER
ARG KUBEBUILDER_VER
# Install kubebuilder
RUN if [ "${INSTALL_KUBEBUILDER}" == "true" ]; then \
        export ARCH=$(case $(uname -m) in x86_64) echo -n amd64 ;; aarch64) echo -n arm64 ;; *) echo -n $(uname -m) ;; esac) && \
        export OS=$(uname | awk '{print tolower($0)}') && \
        export KUBEBUILDER_DL_URL=https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VER} && \
        curl -LO ${KUBEBUILDER_DL_URL}/kubebuilder_${OS}_${ARCH} && \
        chmod +x kubebuilder_${OS}_${ARCH} && \
        mv kubebuilder_${OS}_${ARCH} /usr/local/bin/kubebuilder; \
    fi



# Install k8s code generator
RUN go install k8s.io/code-generator/cmd/client-gen@v0.29.2     && \
    go install k8s.io/code-generator/cmd/deepcopy-gen@v0.29.2   && \
    go install k8s.io/code-generator/cmd/register-gen@v0.29.2   && \
    go install k8s.io/code-generator/cmd/lister-gen@v0.29.2     && \
    go install k8s.io/code-generator/cmd/informer-gen@v0.29.2


# Install required dependencies from go.mod file
# COPY 03-OPERATOR/ /tmp/03-OPERATOR
# RUN cd /tmp/03-OPERATOR/operator && \
#     go mod tidy && \
#     go get ./...