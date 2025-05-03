---
weight: 999
url: "/ansible-a-powerful-agentless-configuration-management-and-orchestrator-solution/"
title: "Ansible: A Powerful Agentless Configuration Management and Orchestrator Solution"
description: "Learn how to use Ansible, an agentless IT automation tool for configuration management, application deployment, and task orchestration"
categories: ["Linux", "Server", "Automation", "DevOps"]
date: "2015-02-26T07:09:00+01:00"
lastmod: "2015-02-26T07:09:00+01:00"
tags: ["Ansible", "Automation", "Configuration Management", "DevOps", "SSH"]
toc: true
---

![Ansible](/images/ansible_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 1.6.2 |
| **Operating System** | Debian 8 |
| **Website** | [Ansible Website](https://www.ansible.com) |
| **Last Update** | 26/02/2015 |
{{< /table >}}

## Introduction

Ansible is an IT automation tool. It can configure systems, deploy software, and orchestrate more advanced IT tasks such as continuous deployments or zero downtime rolling updates.

Ansible's goals are foremost those of simplicity and maximum ease of use. It also has a strong focus on security and reliability, featuring a minimum of moving parts, usage of OpenSSH for transport (with an accelerated socket mode and pull modes as alternatives), and a language that is designed around auditability by humans – even those not familiar with the program.

We believe simplicity is relevant to all sizes of environments and design for busy users of all types – whether this means developers, sysadmins, release engineers, IT managers, and everywhere in between. Ansible is appropriate for managing small setups with a handful of instances as well as enterprise environments with many thousands.

Ansible manages machines in an agentless manner. There is never a question of how to upgrade remote daemons or the problem of not being able to manage systems because daemons are uninstalled. As OpenSSH is one of the most peer reviewed open source components, the security exposure of using the tool is greatly reduced. Ansible is decentralized – it relies on your existing OS credentials to control access to remote machines; if needed it can easily connect with Kerberos, LDAP, and other centralized authentication management systems.[^1]

The documentation is enough complete and well done to avoid rewriting things here. However some tricky things can be done and that's what I'm trying to cover here.

## Installation

To install Ansible, you have 2 choices:

```bash
aptitude install ansible
```

Or with pip (require python-pip):

```bash
pip install ansible
```

## Usage

### Multiple conditions

If you want to set multiple condition, it can sometimes be complicated. Here is an example with 'or' and 'and':

```yaml
- name: remove systemd
  apt: name=systemd state=absent update_cache=yes force=yes
  when: ansible_virtualization_type not in [ 'lxc', 'kvm', 'virtualbox' ] and (ansible_virtualization_role != 'host')
  tags: [common, common_packages, common_services]
```

Here we're do not want to remove systemd package if we're not in a host (ansible_virtualization_role) and if the ansible_virtualization_type variable is not lxc, kvm or virtualbox.

### Debugging vars

You can take a look at a current variable status:

```
- debug: var=kibana_current
```

And see the result:

```yaml
TASK: [kibana | debug var=kibana_current] ************************************* 
<kibana.deimos.lan> ESTABLISH CONNECTION FOR USER: root
ok: [kibana.deimos.lan] => {
    "item": "", 
    "kibana_current": {
        "changed": true, 
        "cmd": "cd /usr/share/nginx/www/kibana ; /usr/bin/git describe --tags ", 
        "delta": "0:00:00.004179", 
        "end": "2014-06-10 12:10:51.797805", 
        "invocation": {
            "module_args": "cd /usr/share/nginx/www/kibana ; /usr/bin/git describe --tags", 
            "module_name": "shell"
        }, 
        "item": "", 
        "rc": 0, 
        "start": "2014-06-10 12:10:51.793626", 
        "stderr": "", 
        "stdout": "v3.1.0", 
        "stdout_lines": [
            "v3.1.0"
        ]
    }
}
```

You can now select a specific item, for example:

```
- debug: var=kibana_current.stdout
```

### Do not notify on changes

To ignore changed=1 for a specific action that will run each time, you can add the change_when statement:

```
- shell: git --git-dir={{kibana_path}}/.git/ describe --tags
  register: kibana_current
  changed_when: false
```

### Force handler to apply

Ansible handlers are defined in a handlers section or file and are called at the end of each play if they have been triggered. This is useful as it means you can have multiple tasks trigger another action, but ensure that the triggered action only runs once:[^2]

```yaml
- name: Add symlink for systemd
  file: src=/lib/systemd/system/mongodb.service dest=/etc/systemd/system/multi-user.target.wants/mongodb.service state=link
  notify: reload systemd

- meta: flush_handlers
 
...
```

This will here force 'reload systemd' handler to be applied without waiting the end of the playbook.

### Check syntax and list tasks

Here is an example to check all your playbook syntax and list tasks at the same time:

```bash
> ansible-playbook -i hosts --syntax-check --list-tasks -e set_env=prod -D --limit server01 site.yml 

playbook: site.yml

  play #1 (physical):
    Ensure iptables is installed (debian)
    Ensure iptables is installed (redhat)
    Prepare iptables rules
    Autoload the rules
    installing docker registry dependencies
    adding jenkins user to docker group
    configure docker
    ensure docker started

  play #2 (common):
    Set hostname to the current machine
    Use Debian CDN in sources.list
    user root
    creating user
    deploy authorized_keys for root
    deploy authorized_keys for users
    generate locales
```

## References

[^1]: http://docs.ansible.com/
[^2]: http://wherenow.org/ansible-handlers/
