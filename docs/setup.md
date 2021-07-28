# Docker-Machine/Rancher 共通事前準備
## 必要なもの
- HPE OneView  
- HPE ProLiantサーバー
- Webサーバー
- (Docker環境を作成する場合) Docker-machineがインストールされた端末
- (Rancher k8s環境を作成する場合) Rancherサーバー

## OSイメージの作成
OSはキックスタートを使用してインストールします。OSイメージはWebサーバーに配置する必要があり、ドライバーからHPE iLO DVD仮想デバイスにマウントされます。  

### 手順例
以下の例ではCentOS 8を使ってキックスタート用のカスタムイメージを作成しています。
OSイメージを作業端末にダウンロード後、マウントさせます。

```
$ mount -t iso9660 -o loop CentOS-8.3.2011-x86_64-minimal.iso /mnt
```

カスタムOSイメージ用のディレクトリを作成して、その中にOSイメージをコピーします。

```
$ mkdir image  
$ find /mnt -maxdepth 1 -mindepth 1 -exec cp -rp {} image/ \;  
```

**grub.cfg**を変更します。キックスタートファイルは後続の手順で*ov-ks*というラベルにするため、下記の例だと*inst.ks=hd:LABEL=ov-ks:/ks.cfg*と設定しています。お好きなラベルを設定して下さい。
*inst.stage2=hd:LABEL* で指定するラベルもお好きなラベルを設定してください。

```
$ vi image/EFI/BOOT/grub.cfg
set default="0" <= 変更

...

set timeout=10 <= 変更

...

### BEGIN /etc/grub.d/10_linux ###
#以下に変更
menuentry 'Install CentOS 8 For Rancher/Docker' --class fedora --class gnu-linux --class gnu --class os {
        linuxefi /images/pxeboot/vmlinuz inst.stage2=hd:LABEL=CentOS-8-3-2011-x86_64-dvd inst.ks=hd:LABEL=ov-ks:/ks.cfg  quiet
        initrdefi /images/pxeboot/initrd.img

```

最後にカスタムOSイメージをisoイメージにします。*-V*オプションで指定するラベルは*inst.stage2=hd:LABEL* で指定したラベルにしてください。

```
$ yum install -y mkisofs
$ cd image 
$ mkisofs \
    -v -r -J -T -l -input-charset utf-8 \
    -o ../CentOS-8.3.2011-x86_64-minimal-ks.iso \
    -b isolinux/isolinux.bin \
    -c isolinux/boot.cat \
    -V 'CentOS-8-3-2011-x86_64-dvd' \
    -no-emul-boot \
    -boot-load-size 4 \
    -boot-info-table \
    -eltorito-alt-boot \
    -e images/efiboot.img \
    -no-emul-boot .
```

## キックスタートファイルの作成
Linux OSインストール用のキックスタートファイルを事前に作成しておく必要があります。キックスタートファイルはWebサーバーに配置する必要があり、ドライバーからHPE iLO Floppy仮想デバイスにマウントされます。また、ネットワークは**インターネットに接続可能**なネットワークを構成してください。 

### 手順例
キックスタートファイルを作成します。キックスタートファイルには自身の環境に合わせたパラメータを記載してください。あくまで下記の記述は例となりますが、以下の要件を全て満たす必要があります。  

- rootユーザーでパスワードを使ってssh接続可能なこと
- rootユーザーで公開鍵を使ってssh接続可能なこと
- インターネット接続可能なネットワークを構成していること


```
# Install OS instead of upgrade
install

# Keyboard layouts
keyboard --vckeymap=us --xlayouts='us'

# Root password
rootpw --plaintext password

# System language
lang en_US.UTF-8

# System authorization information
#auth --enableshadow --passalgo=sha512 <= auth command is not recommend in CentOS8
authselect --useshadow --passalgo sha512

# Run the Setup Agent on first boot
firstboot --enable
ignoredisk --only-use=sda

# SELinux configuration
selinux --disabled

# Firewall configuration
firewall --disabled

# Network information
network  --bootproto=static --device=link --ip=172.16.14.10 --netmask=255.255.240.0 --gateway=172.16.0.1 --nameserver=8.8.8.8 --activate

# Reboot after installation
reboot --eject

# System timezone
timezone Asia/Tokyo

# System bootloader configuration
bootloader --append=" crashkernel=auto" --location=mbr --boot-drive=sda

# Clear the Master Boot Record
zerombr

# Partition clearing information
clearpart --all --initlabel --drive sda

# Disk partitioning information
part /boot/efi --fstype="efi" --ondisk=sda --size=1024 --label=bootef --fsoptions="umask=0077,shortname=winnt"
part /boot --fstype="xfs" --ondisk=sda --size=1024 --label=boot
part pv.186 --fstype="lvmpv" --ondisk=sda --size=1 --grow
volgroup centos --pesize=4096 pv.186
logvol / --fstype="xfs" --size=1 --grow --name=root --vgname=centos --label=root

# Others
eula --agreed
xconfig  --startxonboot

%packages
%end

%post
# hosts
echo "172.16.14.10  test01.hybrid-lab.local  test01" >> /etc/hosts

# sshd accept root login
sed -i 's/^#PermitRootLogin yes/PermitRootLogin yes/' /etc/ssh/sshd_config
sed -i 's/^#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config

# fixed OS and Kernel version in the case of Cenos8
echo "excludepkgs=centos*,kernel*" >> /etc/dnf/dnf.conf
# fixed OS and Kernel version in the case of Cenos7
echo "excludepkgs=centos*,kernel*" >> /etc/yum.conf

%end
```

キックスタートファイル作成後、isoイメージ変換します。その際、iso-imageのラベルはカスタムOSイメージを作成した際に使用したラベル名を設定します。下記の例では*ov-ks*を指定しています。  
**isoイメージの名前は"IPアドレス.iso"としてください。**ドライバーはオプションの引数として与えられたサーバーのIPアドレスを元にキックスタートファイルを探します。

```
$ ls ks
ks.cfg
$ mkisofs -r -J -V ov-ks -o 172.16.14.10.iso ks
```

## Webサーバーへ各種ファイルを配置
先の手順で作成した各種ファイルをwebサーバーに配置します。

- カスタムOSイメージのiso
- キックスタートファイルのiso

配置する場所はwebサーバー経由でダウンロードできる場所ならどこでも構いません。  
配置後、作業端末から各種ファイルがダウンロードできるかを確認してください。

## HPE OneViewでサーバープロファイルテンプレートを用意する
サーバープロファイルテンプレートはDocker/Rancher k8sを構成する際にドライバーが使用します。
ターゲットとなるサーバーとDocker machine/Rancherが疎通可能なネットワークをサーバープロファイルテンプレートに設定しておく必要があります。  
また、OSをインストールするためのブートディスクの構成も適切にサーバープロファイルテンプレートで定義しておく必要があります。
