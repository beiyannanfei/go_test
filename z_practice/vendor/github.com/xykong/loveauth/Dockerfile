FROM centos:7
MAINTAINER by xy.kong@gmail.com

# install openssh-server
RUN yum -y install openssh-server openssl passwd git zsh go wget

RUN ssh-keygen -q -t rsa -b 2048 -f /etc/ssh/ssh_host_rsa_key -N ''
RUN ssh-keygen -q -t ecdsa -f /etc/ssh/ssh_host_ecdsa_key -N ''
RUN ssh-keygen -t dsa -f /etc/ssh/ssh_host_ed25519_key  -N ''

RUN sed -i "s/#UsePrivilegeSeparation.*/UsePrivilegeSeparation no/g" /etc/ssh/sshd_config
RUN sed -i "s/UsePAM.*/UsePAM no/g" /etc/ssh/sshd_config

RUN /bin/echo -e "LANG=\"en_US.UTF-8\"" > /etc/default/local
RUN mkdir -p ~/.ssh
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDOd6C+0kpbdihkS6BG7YIkpaxwJysjQbGLs4sub9aUGY05otDul8mzQcg2cdVBo6tOzmmAv8N8bA1y6NsbPoqibMlC5c+I25Vx0g2bINsbQR8ViX5haYekKjy9EYafdM7Q9xeMi3SYjW5BprfdUQL/5tSPMX5Umfr7PhAwm4bHFOQl1/OT3kUulhiR9sFmux/tvAGwUasVZQBJKDgsHCug3aG4CTdQMcMFsglCGuv/77bK7YhSwqdoJPPnuWN+t9p7955C4bO+j4uW3BGq6ID7CvzeTRzTFBLWhhmM9f19NLBzqfMW5NG/+AfNR5zWfVb4NPhd42qXPHmGr9jthswNimnT06oJfyMODcyPZFsMsKAsbTq1BzLTbQ/dY8kXP6Ck1S3kJpjuRE0MFgj+w5NWdeb9R6bSwpR1YCE/zaGLTqB392IpjANra68A+osp5F8q7XaLHhkWyFqW2GjMAtdw9S7gXAMchjui4UZTjnSXCUz+kDxTAFeryKWkQaQhDNd/CSAYrrUCQogIQ8TtYHkpH27GADVlP2sNnG9sixJIcIBIW7KyvWy+vhstn9Nvf8Y0USydqU7MXKNiGx9U3o6uRKoAsFlO0iTq9cPDtg5yHZs/EDrHoRXOqHD+PBeyD6bZf2mlc8hDEp10RD49ebkQ4qCUKt1DCde4RaakvC9+2w== xy.kong@gmail.com" > ~/.ssh/authorized_keys


# install zsh

RUN sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
RUN echo $'local ret_status="%(?:%{$fg_bold[green]%}➜ :%{$fg_bold[red]%}➜ )"\n\
PROMPT="%(!.%{%F{yellow}%}.)$USER @ %{$fg[white]%}%M %{$fg_bold[red]%}➜ %{$fg_bold[green]%}%p %{$fg[cyan]%}%c %{$fg_bold[blue]%}$(git_prompt_info)%{$fg_bold[blue]%} % %{$reset_color%}"\n\
ZSH_THEME_GIT_PROMPT_PREFIX="%{$fg_bold[blue]%}git:(%{$fg[red]%}"\n\
ZSH_THEME_GIT_PROMPT_SUFFIX="%{$reset_color%} "\n\
ZSH_THEME_GIT_PROMPT_DIRTY="%{$fg[blue]%}) %{$fg[yellow]%}✗"\n\
ZSH_THEME_GIT_PROMPT_CLEAN="%{$fg[blue]%})"\n\
' > ~/.oh-my-zsh/custom/themes/robbyrussell.zsh-theme


# install mysql

RUN wget http://dev.mysql.com/get/mysql-community-release-el7-5.noarch.rpm
RUN rpm -ivh mysql-community-release-el7-5.noarch.rpm
RUN yum install -y mysql-community-server
RUN mysql_install_db --user=mysql --ldata=/var/lib/mysql


# install redis
RUN yum install -y epel-release
RUN yum install -y redis


# install loveauth

RUN ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN printf $'[user]\n\
    email = xy.kong@gmail.com\n\
    name = xykong\n\n\
[url "git@github.com:"]\n\
	insteadOf = https://github.com/\n\
'> ~/.gitconfig


#RUN go get -u -v github.com/xykong/loveauth

#    go get -u -v github.com/xykong/loveauth;\
ADD ~/.ssh/xy.kong@gmail.com.key /tmp/
RUN ssh-agent /tmp
RUN go get -u -v github.com/xykong/loveauth
RUN rm /tmp/mykey


EXPOSE 22

#CMD /usr/sbin/sshd -D

CMD /usr/bin/mysqld_safe &\
    /usr/bin/redis-server &\
    /usr/sbin/sshd -D


