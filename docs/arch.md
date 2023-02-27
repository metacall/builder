## Arch

The document describes the architecture of the project.

### Image Dependencies

According to the [Dockerfiles](https://github.com/metacall/core/tree/develop/tools) in core,
we infer the dependencies and inheritance relationships between different versions of images.

This determines the dependencies and stages of Builder subcommands.

```mermaid
graph LR
    Base(Base Image)
    Runtime(Runtime Image)
    Deps(Deps Image)
    Dev(Dev Image)
    CLI(CLI Image)
    
    subgraph Image Dependencies
        direction LR
        Base --> Runtime
        Base --> Deps
        Deps --> Dev
        Dev --> Runtime
        Dev --> CLI
        Runtime --> CLI
    end
```

### Builder Arch

According to the [Dockerfiles](https://github.com/metacall/core/tree/develop/tools)
and [Shell Scripts](https://github.com/metacall/core/tree/develop/tools) in core,
Builder is designed as a multi-layer architecture, with each layer providing capabilities to the upper layers
and enables the Builder to flexibly merge and diff image layers.

```mermaid
flowchart TB
    CLI-Builder(Builder)
    CLI-Buildctl(Buildctl)
    CLI-Registry(Registry)
    
    Flags-Options(Options)
    Flags-Languages(Languages)
    
    Staging-Dev(Dev)
    Staging-Runtime(Runtime)
    Staging-Deps(Deps)
    Staging-CLI(CLI)
    
    Builder-Environment(Environment)
    Builder-Build(Build)
    Builder-Configure(Configure)
    Builder-More(......)
    
    Buildkit-LLB(LLB)
    Buildkit-Op(Merge&Diff)
    
    Comment-CLI(End User)
    Comment-Flags(Customizable & Minimal)
    Comment-Staging(High-Level Image)
    Comment-Builder(Specific Instructions)
    Comment-Buildkit(Image Capabilities)
    
    subgraph CLI
        direction LR
        CLI-Builder --- CLI-Buildctl --- CLI-Registry
        
    end
    
    subgraph Flags
        direction LR
        Flags-Options --- Flags-Languages
    end
    
    subgraph Staging
        direction LR
        Staging-Dev --- Staging-Runtime --- Staging-Deps --- Staging-CLI
    end
    
    subgraph Builder
        direction LR
        Builder-Environment --- Builder-Build --- Builder-Configure --- Builder-More
    end
    
    subgraph Buildkit
        direction LR
        Buildkit-LLB --- Buildkit-Op
    end
    
    %% To make above links invisible. But GitHub don't support '~~~' syntax. 
    linkStyle 0,1,2,3,4,5,6,7,8,9 display:none
    
    subgraph Builder Arch
        direction LR
        CLI --> Comment-CLI
        Flags --> Comment-Flags
        Staging --> Comment-Staging
        Builder --> Comment-Builder
        Buildkit --> Comment-Buildkit
    end
```
