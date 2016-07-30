Senter controller app
====================

git clone https://github.com/snub/senter.git
cd senter
godep restore
go get github.com/jinzhu/inflection

git clone https://github.com/snub/senter_ctrl.git
cd senter_ctrl
godep restore
go build
