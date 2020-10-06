import datetime
import os

if __name__ == '__main__':
    current_path = os.getcwd()
    project_path = os.path.abspath(os.path.join(current_path, ".."))
    name = "jdsl"
    package = "shared"
    now = datetime.datetime.now()
    str_now = now.strftime("%Y-%m-%d-%H-%M-%S")
    out = os.path.join(current_path, 'xgo', name + '_' + str_now)
    go_cache = os.getenv("GOCACHE")
    go_path = os.getenv("GOPATH")
    mod_path = os.path.join(go_path, "pkg", "mod")
    args = {
        "out": out,
        "project_path": project_path,
        "mod_path": mod_path,
        "go_cache": go_cache,
        "name": name,
        "package": package
    }
    a = 'Docker run ' \
        '--rm ' \
        '-v {out}:/build ' \
        '-v {project_path}:/source ' \
        '-v {mod_path}:/go/pkg/mod ' \
        '-v {go_cache}:/go-build ' \
        '-e OUT={name} ' \
        '-e PACK={package} ' \
        '-e FLAG_V=false ' \
        '-e FLAG_X=false ' \
        '-e FLAG_RACE=false ' \
        '-e FLAG_LDFLAGS="-s -w" ' \
        '-e FLAG_BUILDMODE=c-shared ' \
        '-e TARGETS=./. ' \
        '-e GO111MODULE=on ' \
        '-e GOPROXY=https://goproxy.io ' \
        'crazymax/xgo ' \
        '.'
    cmd = a.format(
        out=out,
        project_path=project_path,
        mod_path=mod_path,
        go_cache=go_cache,
        name=name,
        package=package
    )
    print (cmd)
    os.system(cmd)
