[![Build Status](https://travis-ci.com/HPE-Japan-Presales/docker-machine-driver-ov.svg?branch=main)](https://travis-ci.com/HPE-Japan-Presales/docker-machine-driver-ov)

[English](/README_en.md)
# Docker-Machine/Rancher HPE Servers Driver
このDocker-Machine/RancherドライバーはHPE OneViewで管理されたHPE Serversに向けて作成されています。  
Serverハードウェアの設定からOSのインストール、DockerまたはRancher Kubernetesのインストールまでを自動化します。  
Hewlett Packardのレポジトリには既に[docker-machine-oneview](https://github.com/HewlettPackard/docker-machine-oneview)がありますが、メンテナンスを終了している様子なのでこのドライバーを作成しました。

## テスト環境

### Docker-machine

|  サーバーモデル | OS |  OneViewバージョン  |  Docker-Machineバージョン  |  その他  |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- |
|   HPE Synergy480 Gen9  |  CentOS 7.8  |   5.30.00  |   0.16.2  |    |


### Rancher

|  サーバーモデル | OS |  OneViewバージョン  |  Rancherバージョン  |  Rancher k8sバージョン  |  その他  |
| ---- | ---- | ---- | ---- | ---- | ---- | ---- |
|   HPE Synergy480 Gen9  |  CentOS 7.8  |   5.30.00  |   2.5.7      |  v1.17.17  |      | 
|   HPE Synergy480 Gen9  |  CentOS 7.8  |   5.30.00  |   2.5.9-rc4  |  v1.17.17  |      | 
|   HPE Synergy480 Gen9  |  CentOS 7.8  |   5.30.00  |   2.5.9-rc4  |  v1.18.20  |      | 
|   HPE Synergy480 Gen9  |  CentOS 7.8  |   5.30.00  |   2.5.9-rc4  |  v1.19.12  |      |



## 各種手順
- [事前準備](docs/setup.md)
- [Docker-machineでの使用方法](docs/docker-machine/setup.md)
- [Rancherでの使用方法](docs/rancher/setup.md)

## アーキテクチャ
### 概要
```
                           │
                           │
┌─────────────────────┐    │    ┌─────────────────────┐
│HPE OneView          │    │    │Web Server           │
│                     │    │    │  To Provide         │
│                     ├────┼────┤   OS Image&Kickstart│
│                     │    │    │                     │
└─────────────────────┘    │    └─────────────────────┘
                           │
┌─────────────────────┐    │    ┌─────────────────────┐
│HPE Server           │    │    │Docker Machine       │
│  Managed By OneView ├────┼────┤ Or Rancher          │
│                     │    │    │                     │
│                     │    │    │                     │
└─────────────────────┘    │    └─────────────────────┘
```
ドライバーを動作させるためには4つのコンポーネントが必要となります。  
- HPEサーバー(Docker/Rancher k8s構築ターゲット)  
- Docker MachineまたはRancher環境  
- HPE OneView  
- Webサーバー  

構築対象となるサーバーはHPE OneViewで管理されたHPEサーバーのみが使用できます。

### OS自動インストール方法
```
┌──────────────────────┐           ┌──────────────────────┐
│      HPE Server      │           │       Web Server     │
│                      │           │                      │
│   iLO Virtual Mount  │           │    - OS Image        │
│    - OS Image        ├───────────┤    - Kick Start File │
│    - Kick Start File │           │                      │
│                      │           │                      │
│                      │           │                      │
└──────────────────────┘           └──────────────────────┘
```
OS自動インストール仕組みはHPE iLO Virtual Mountを利用しています。  
事前に用意したWeb Server上にベースとなるOSイメージとキックスタートファイルを配置します。
OSのベースイメージは仮想DVDデバイスとしてマウントされ、キックスタートファイルは仮想Floppyデバイスとしてマウントされます。事前に用意したOSベースイメージは仮想Floppyデバイスをキックスタートファイルとして認識させるためにカスタマイズが必要です。  



## その他
- [Known Issue](docs/known_issue.md)