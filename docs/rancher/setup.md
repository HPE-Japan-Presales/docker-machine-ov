# Rancherでの使用方法
Rancherでドライバーを使用するための手順を記載しています。  
事前に[こちらの準備](../setup.md)が終了していることを確認してください。

## ノードドライバーのセットアップ
ノードドライバーをRancherにインストールします。  
RahcnerからTools > Driversを選択します。
![](add_node_driver01.png "")

Node DriversタブからAdd Node Driverを選択します。
![](add_node_driver02.png "")

Download URLにGit hub releaseのURLを設定します。または環境内のWebサーバーにドライバーをダウンロードして公開した後、そのURLを設定します。
![](add_node_driver03.png "")

ドライバーのインストールが完了するとドライバーがActiveになります。
![](add_node_driver04.png "")

## クラスタの作成
Rancher k8sクラスタを作成します。Add Clusterボタンからクラスタ作成メニューに移動します。
![](add_cluster.png "")

ノードドライバーは先ほど追加した"Ov"を選択します。
![](node_driver_menu.png "")

クラスターセットアップのパラメータを入力します。Cluster Name, Name Prefix、役割にチェックした後、ノードドライバー用のパラメータテンプレートを作成します。
![](cluster_setup01.png "")

ノードドライバーのパラメータを設定します。
![](node_driver_params.png "")

ノードドライバー用のパラメータテンプレート作成後、Rancher k8sの各種設定を選択します。設定後、Createボタンからクラスタを作成します。
![](cluster_setup02.png "")

プロビジョニングが開始します。HPE OneViewnでサーバープロファイルの作成が開始されます。
![](provisioning.png "")

30-40分(環境により異なる)経つとRancher k8sクラスタが完成します。
![](clusters.png "")

ダッシュボードからクラスタの状態が見えると思います。
![](dashboard.png "")