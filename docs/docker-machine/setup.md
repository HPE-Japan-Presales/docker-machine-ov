# Docker-machineでの使用方法
Docker-machineでドライバーを使用するための手順を記載しています。  
事前に[こちらの準備](../setup.md)が終了していることを確認してください。

## 作業端末にドライバーをインストールする
Docker-machineコマンドがインストールされている作業端末にドライバーをインストールします。
任意のバージョンのドライバーを[ダウンロード](https://github.com/HPE-Japan-Presales/docker-machine-ov/releases)してください。

```
$ tar xvf docker-machine-driver-ov-v0.0.10-linux-amd64.tar.gz
$ cp -p docker-machine-driver-ov /usr/local/bin
$ docker-machine create --driver ov --help
...
   --ov-debug												(Option) Debug Flag For This Driver [$OV_DEBUG]
   --ov-oneview-api-version "1800"									HPE OneView: OneView API Version [$OV_ONEVIEW_API_VERSION]
   --ov-oneview-domain 											HPE OneView: (Option) OneView Domain [$OV_ONEVIEW_DOMAIN]
   --ov-oneview-endpoint "https://oneview.hpe.com"							HPE OneView: OneView Endpoint Address [$OV_ONEVIEW_ENDPOINT]
   --ov-oneview-password "password"									HPE OneView: OneView User Password [$OV_ONEVIEW_PASSWORD]
   --ov-oneview-server-hardware 									HPE OneView: Target Server Hardware Name In OneView(Exactly Same As OneView Displayed, Need Spaces Between Strings) [$OV_ONEVIEW_SERVER_HARDWARE]
   --ov-oneview-server-profile-template 								HPE OneView: OneView Server Profile Template Name For Target Server Hardware [$OV_ONEVIEW_SERVER_PROFILE_TEMPLATE]
   --ov-oneview-user "administrator"									HPE OneView: OneView User [$OV_ONEVIEW_USER]
   --ov-server-address 											New Server: Target Server Address [$OV_SERVER_ADDRESS]
   --ov-server-kickstart-base-url 									New Server: Kickstart Image Base URL. If your kickstart is on http://web01/docker/kickstart.iso, you shoud set this value as http://web01/docker. [$OV_SERVER_KICKSTART_BASE_URL]
   --ov-server-os-url 											New Server: OS Image URL [$OV_SERVER_OS_URL]
   --ov-server-root-password "password"									New Server: Target Server Root Password [$OV_SERVER_ROOT_PASSWORD]
   --ov-yaml 	
...
```

## Docker環境の構築
YAMLファイル、またはコマンド引数で各種パラメータを設定することでDocker-machineからDocker環境を作成できます。  
以下はYAMLファイルにパラメータを記載した例です。

```
oneview:
  endpoint: "https://<YOUR ONEVIEW IP>"
  api-version: 1200
  user: "<YOUR ONEVIEW USER>"
  password: "<YOUR ONEVIEW USER PASSWORD>"
  domain: ""
  server-profile-template: "Rancher-template"
  server-hardware: "SGH652SV73, bay 5"
server:
  address: "<YOUR SERVER IP>"
  root-password: "password"
  kickstart-base-url: "http://<YOUR WEB SERVER IP>/tak"
  os-url: "http://<YOUR WEB SERVER IP>/tak/CentOS-7-x86_64-Minimal-2003-ks.iso"
```

以下のYAMLを指定してdocker-machineを起動させます。

```
$ sudo docker-machine create --driver ov --ov-yaml ./configs/examples/synergy.yaml test01
(test01) Configuration is read from yaml
Running pre-create checks...
(test01) Check HPE OneView configurations
(test01) Check new server configurations
...
<About 40min>
...
$ docker-machine-driver-ov git:(dev) ✗ docker-machine ls                                                                 
NAME     ACTIVE   DRIVER   STATE     URL                       SWARM   DOCKER     ERRORS
test01   -        ov       Running   tcp://172.16.14.10:2376           v20.10.7                                                              

```

## コマンドオプション
| コマンドオプション名 | 環境変数 | YAML | 型 | デフォルト値 | 説明 |
| ------------- | ------------- | ------------- | ------------- | ------------- | ------------- |
| --ov-yaml  | OV\_YAML  | N/A  | string  | None  | YAMLファイルのパスを指定します。YAMLファイルを指定した場合はその他のオプションは必要ありません。  |
| --ov-oneview-endpoint  | OV\_ONEVIEW\_ENDPOINT  | oneview.endpoint  | string  |None  | HPE OneViewのエンドポイントを指定します。</br> (例 http://oneview.hpe.com) |
| --ov-oneview-api-version  | OV\_ONEVIEW\_API\_VERSION  | oneview.api-version  | int  | 1800  | HPE OneView APIバージョンを指定してます。  |
| --ov-oneview-user  | OV\_ONEVIEW\_USER  |  oneview.user   | string  |  administrator  | HPE OneViewのユーザー名を指定します。ユーザーはインフラ管理者以上の権限を持っている必要があります。  |
| --ov-oneview-password  | OV\_ONEVIEW\_PASSWORD  | oneview.password  | string  |  password | HPE OneViewのユーザーパスワードを指定します。  |
| --ov-oneview-domain  | OV\_ONEVIEW\_DOMAIN  | oneview.domain  | string  | None  | (オプション) HPE OneViewドメイン名を指定します。  |
| --ov-oneview-server-profile-template  | OV\_ONEVIEW\_SERVER\_PROFILE\_TEMPLATE  | oneview.server-profile-template  | string  | None  | HPE OneView上に作成されたサーバープロファイルテンプレート名を指定します。このテンプレートはサーバー作成の際に使用されます。  |
| --ov-oneview-server-hardware  | OV\_ONEVIEW\_SERVER\_HARDWARE  | oneview.server-hardware  | string  | None  | HPE OneView上に登録されたサーバーハードウェア名を指定します。このサーバーは実際にDocker/Rancher k8sが作成される対象のサーバーとなります。  |
| --ov-server-address  | OV\_SERVER\_ADDRESS  | server.address  | string   | None  | 作成するサーバーのIPアドレスを指定します。IPアドレスは事前準備したキックスタートファイル内に定義されたIPアドレスです。 |
| --ov-server-root-password | OV\_SERVER\_ROOT\_PASSWORD  | server.root-password  | string   | password  | 作成するサーバーのRootパスワードを指定します。Rootパスワードは事前準備したキックスタートファイル内に定義されたRootパスワードです。  |
| --ov-server-kickstart-base-url  | OV\_SERVER\_KICKSTART\_BASE\_URL  | server.kickstart-url  | string   | None  | キックスターファイルイメージのベースURLを指定します。<br>(例: もしhttp://web-server/rancher/172.16.1.10.iso というURLにキックスタートファイルがある場合、http://web-server/rancher を指定してください。)  |
| --ov-server-image-url  | OV\_SERVER\_IMAGE\_URL  | server.image-url  | string   | None  | OSイメージのURLを指定します。</br>(例：http://webserver/rancher/centos7.iso) |
| --ov-debug  | OV\_DEBUG  | N/A  | string  | None  | (オプション)デバッグの際に指定してください。  |
