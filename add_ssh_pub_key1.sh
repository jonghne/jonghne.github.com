#!/usr/bin/env bash
set -v
key=`cat ~/.ssh/id_rsa.pub`

dosth="mkdir -p ~/.ssh && echo '$key' >> ~/.ssh/authorized_keys"

auto_ssh_c () {
    expect -c " set timeout -1;
                spawn ssh $1@$2 $dosth;
                expect {
                    *(yes/no)* {send -- yes\r;exp_continue;}
                    *assword:* {send -- $3\r;
                                 expect {
                                    *denied* {exit 2;}
                                    eof
                                 }
                    }
                    eof         {exit 1;}
                }
                "
    return $?
}

auto_ssh_c $1 $2 $3
