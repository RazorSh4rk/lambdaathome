export type RuntimeConfig = {
    name: string;
    tag: string;
    runtime: string;
    port: string;
    volume: string;
    source: string;
    id: string;
};

interface Port {
    IP: string;
    PrivatePort: number;
    PublicPort: number;
    Type: string;
}

interface Network {
    IPAMConfig: null;
    Links: null;
    Aliases: null;
    MacAddress: string;
    DriverOpts: null;
    NetworkID: string;
    EndpointID: string;
    Gateway: string;
    IPAddress: string;
    IPPrefixLen: number;
    IPv6Gateway: string;
    GlobalIPv6Address: string;
    GlobalIPv6PrefixLen: number;
    DNSNames: null;
}

interface Mount {
    Type: string;
    Source: string;
    Destination: string;
    Mode: string;
    RW: boolean;
    Propagation: string;
}

interface HostConfig {
    NetworkMode: string;
}

interface NetworkSettings {
    Networks: {
        [key: string]: Network;
    };
}

export type Container = {
    Id: string;
    Names: string[];
    Image: string;
    ImageID: string;
    Command: string;
    Created: number;
    Ports: Port[];
    Labels: {};
    State: string;
    Status: string;
    HostConfig: HostConfig;
    NetworkSettings: NetworkSettings;
    Mounts: Mount[];
}