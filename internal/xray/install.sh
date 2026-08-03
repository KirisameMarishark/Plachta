xray_install() {


    mkdir -p "${ROOT_DIR}/runtime/xray"

    unzip -o \
        "${ROOT_DIR}/runtime/Xray-linux-64.zip" \
        -d "${ROOT_DIR}/runtime/xray"

    install -Dm755 \
        "${ROOT_DIR}/runtime/xray/xray" \
        /usr/local/bin/xray

    log_info "Xray installed to /usr/local/bin/xray"

}