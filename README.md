#  dotf - A hugo-like solution to dot files


## Usage


### Config

Config directory consists of templates and config file that maps tpl to destination path

**Example Config**
```YAML
- name: "bash" # Profile name (can be anything)
  templates:
    - bash/*.tpl # list of glob template paths to render
  filemap: # map of template paths to their destination files, leading directories will be created
    bashrc.tpl: ".bashrc" # maps bashrc.tpl to $BASE/.bashrc
    bash_profile.tpl: ".bash_profile"
```
**Example Template (./bash/.bashrc)**
```
export HOME="{{.var.home }}"
export GPG_KEYID="{{ .var.email }}"

export GOOS="{{.runtime.GOOS}}"
export GOARCH="{{.runtime.GOARCH}}"

```

## Examples


### **Local file** 

`dot install --base /home/testuser ./mydots/  ./my_custom_vars.yml
`
### **Remote Config** 

`dot install git+https://github.com/rtgnx/example-dots/  ./my_custom_vars.yml`

### **Remote Config and Remote Variables** 

`dot install git+https://github.com/rtgnx/example-dots/  git+https://github.com/me/custom-vars`

### **Remote with specific branch or tag**
```
 dot install --base /home/testuser git+https://github.com/rtgnx/example-dots#refs/tags/v0.1.0 \
   git+https://github.com/me/custom-vars#refs/heads/develop
```

[Supported Schemes](https://github.com/hairyhenderson/go-fsimpl)