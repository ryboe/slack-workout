# -*- mode: ruby -*-

require 'fileutils'

Vagrant.require_version ">= 1.7.4"

PROJECT_PATH = File.dirname(__FILE__)
CLOUD_CONFIG_PATH = File.join(PROJECT_PATH, "secret/user-data")

Vagrant.configure("2") do |config|
    config.ssh.insert_key = false  # always use Vagrant's insecure default key
    config.vm.box = "coreos-stable"
    config.vm.box_url = "http://stable.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json"
    vm_name = "slack-pushups-host"
    config.vm.hostname = vm_name

    # will appear as "slack-pushups-host" in Vagrant CLI
    config.vm.define vm_name do |name|
    end

    config.vm.provider :virtualbox do |vb|
        # CoreOS doesn't support VirtualBox guest additions, which are device
        # and system apps that integrate the guest and host operating systems
        # with things like mouse controls, accelerated graphics, and a shared
        # clipboard.
        vb.check_guest_additions = false
        vb.functional_vboxsf = false  # not supported
        vb.memory = 1024
        vb.cpus = 1
        vb.name = vm_name
    end

    # set box static ip, not available to public internet (through network bridge)
    config.vm.network :private_network, ip: "172.16.0.2"

    # copy cloud-config into box
    config.vm.provision :file, :source => CLOUD_CONFIG_PATH, :destination => "/tmp/vagrantfile-user-data"
    config.vm.provision :shell, :inline => "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant", :privileged => true
end
