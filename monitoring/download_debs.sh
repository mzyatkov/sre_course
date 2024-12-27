#!/usr/bin/env bash
set -euo pipefail

# Задаем версии компонентов
PROMETHEUS_VERSION="2.45.0"
ALERTMANAGER_VERSION="0.25.0"
BLACKBOX_EXPORTER_VERSION="0.23.0"
NODE_EXPORTER_VERSION="1.6.1"
POSTGRES_EXPORTER_VERSION="0.16.0"
PATRONI_EXPORTER_VERSION="0.3.0"   
ETCD_EXPORTER_VERSION="0.10.0"    

# Создаем нужные директории
mkdir -p roles/prometheus/files
mkdir -p roles/node_exporter/files
mkdir -p roles/postgres_exporter/files
mkdir -p roles/patroni_exporter/files
mkdir -p roles/etcd_exporter/files

# Рабочая директория для временных файлов
TMP_DIR="$(mktemp -d)"

########################################
# Prometheus
########################################
PROMETHEUS_URL="https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.linux-amd64.tar.gz"
echo "Downloading Prometheus..."
wget -O "${TMP_DIR}/prometheus.tar.gz" "${PROMETHEUS_URL}"
tar -xzf "${TMP_DIR}/prometheus.tar.gz" -C "${TMP_DIR}"
cp -rf "${TMP_DIR}/prometheus-${PROMETHEUS_VERSION}.linux-amd64/prometheus" roles/prometheus/files/
cp -rf "${TMP_DIR}/prometheus-${PROMETHEUS_VERSION}.linux-amd64/promtool" roles/prometheus/files/

########################################
# Alertmanager
########################################
ALERTMANAGER_URL="https://github.com/prometheus/alertmanager/releases/download/v${ALERTMANAGER_VERSION}/alertmanager-${ALERTMANAGER_VERSION}.linux-amd64.tar.gz"
echo "Downloading Alertmanager..."
wget -O "${TMP_DIR}/alertmanager.tar.gz" "${ALERTMANAGER_URL}"
tar -xzf "${TMP_DIR}/alertmanager.tar.gz" -C "${TMP_DIR}"
cp -rf "${TMP_DIR}/alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/alertmanager" roles/prometheus/files/
cp -rf "${TMP_DIR}/alertmanager-${ALERTMANAGER_VERSION}.linux-amd64/amtool" roles/prometheus/files/

########################################
# Blackbox Exporter
########################################
BLACKBOX_URL="https://github.com/prometheus/blackbox_exporter/releases/download/v${BLACKBOX_EXPORTER_VERSION}/blackbox_exporter-${BLACKBOX_EXPORTER_VERSION}.linux-amd64.tar.gz"
echo "Downloading Blackbox Exporter..."
wget -O "${TMP_DIR}/blackbox_exporter.tar.gz" "${BLACKBOX_URL}"
tar -xzf "${TMP_DIR}/blackbox_exporter.tar.gz" -C "${TMP_DIR}"
cp -rf "${TMP_DIR}/blackbox_exporter-${BLACKBOX_EXPORTER_VERSION}.linux-amd64/blackbox_exporter" roles/prometheus/files/

########################################
# Node Exporter
########################################
NODE_URL="https://github.com/prometheus/node_exporter/releases/download/v${NODE_EXPORTER_VERSION}/node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64.tar.gz"
echo "Downloading Node Exporter..."
wget -O "${TMP_DIR}/node_exporter.tar.gz" "${NODE_URL}"
tar -xzf "${TMP_DIR}/node_exporter.tar.gz" -C "${TMP_DIR}"
cp -rf "${TMP_DIR}/node_exporter-${NODE_EXPORTER_VERSION}.linux-amd64/node_exporter" roles/node_exporter/files/

########################################
# Postgres Exporter (postgres_exporter)
########################################
# Официального Postgres Exporter от Prometheus нет, есть проект: https://github.com/prometheus-community/postgres_exporter
# Для примера возьмём бинарник оттуда (если есть готовый релиз)
POSTGRES_EXPORTER_URL="https://github.com/prometheus-community/postgres_exporter/releases/download/v${POSTGRES_EXPORTER_VERSION}/postgres_exporter-${POSTGRES_EXPORTER_VERSION}.linux-amd64.tar.gz"
echo "Downloading Postgres Exporter..."
wget -O "${TMP_DIR}/postgres_exporter.tar.gz" "${POSTGRES_EXPORTER_URL}"
tar -xzf "${TMP_DIR}/postgres_exporter.tar.gz" -C "${TMP_DIR}"
cp -rf "${TMP_DIR}/postgres_exporter-${POSTGRES_EXPORTER_VERSION}.linux-amd64/postgres_exporter" roles/postgres_exporter/files/

########################################
# Patroni Exporter
########################################
# Patroni Exporter - Showmax/patroni-exporter (https://github.com/Showmax/patroni-exporter)
# Допустим, есть сборка:
PATRONI_EXPORTER_URL="https://github.com/ccakes/patroni_exporter/releases/download/v${PATRONI_EXPORTER_VERSION}/patroni_exporter_linux_amd64"
echo "Downloading Patroni Exporter..."
wget -O "roles/patroni_exporter/files/patroni_exporter" "${PATRONI_EXPORTER_URL}"

echo "All downloads and extractions completed."

# Очистка временного каталога
rm -rf "${TMP_DIR}"
