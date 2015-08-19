# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.box = "https://cloud-images.ubuntu.com/vagrant/utopic/current/utopic-server-cloudimg-amd64-vagrant-disk1.box"

  config.vm.synced_folder "/Users/#{ENV['USER']}/work/go", "/var/go"
  config.vm.provision :shell, path: "bootstrap.sh"
end
