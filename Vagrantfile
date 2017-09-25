VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |cfg|
    cfg.vm.box = "ubuntu/xenial64" # 16.04 LTS
    cfg.vm.hostname = "scrapper-service"
    cfg.vm.provision :shell, path: "provision.sh"
    cfg.vm.synced_folder ".", "/go/src/github.com/dk13danger/scrapper-service"
    cfg.vm.provider :virtualbox do |v|
        v.memory = "2048"
    end
end
