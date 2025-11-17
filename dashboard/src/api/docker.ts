import { post, type ApiResponse } from "./base";

export interface DockerInfo {
    version: string;
    api_version: string;
    os: string;
    arch: string;
    total_memory: number;
    total_cpu: number;
    total_disk: number;
}

export interface Container {
    id: string;
    name: string;
    image: string;
    status: string;
    created: number;
    ports: string[];
}

export interface Image {
    Id: string;
    Size: number;
    Created: number;
    Containers: number;
    Labels?: Record<string, string>;
    ParentId?: string;
    RepoDigests?: string[];
    RepoTags: string[];
}

export interface Volume {
    Name: string;
    CreatedAt: string;
    Size: number;
    Driver: string;
    Labels?: Record<string, string>;
    Mountpoint: string;
    Scope: string;
}


export interface Network {
    Id: string;
    Name: string;
    Created: string;
    Scope: string;
    Driver: string;
    EnableIPv4: boolean;
    EnableIPv6: boolean;
    IPAM: {
        Driver: string;
        Config?: {
            Subnet?: string;
            Gateway?: string;
        }[];
    };
    Internal: boolean;
    Attachable: boolean;
    Ingress: boolean;
    Containers: {
        Id: string;
        Name: string;
    }[];
}

const getHeaders = (nodeId: number) => {
    return {
        "X-Node-ID": nodeId.toString(),
    };
};

export const dockerApi = {
    async info(nodeId: number): Promise<ApiResponse<DockerInfo>> {
        const header = getHeaders(nodeId);
        return post<DockerInfo>("/api/docker/info", null, header);
    },

    async containers(nodeId: number): Promise<ApiResponse<Container[]>> {
        const header = getHeaders(nodeId);
        return post<Container[]>("/api/docker/containers", null, header);
    },

    async images(nodeId: number): Promise<ApiResponse<Image[]>> {
        const header = getHeaders(nodeId);
        return post<Image[]>("/api/docker/images", null, header);
    },

    async volumes(nodeId: number): Promise<ApiResponse<Volume[]>> {
        const header = getHeaders(nodeId);
        return post<Volume[]>("/api/docker/volumes", null, header);
    },

    async networks(nodeId: number): Promise<ApiResponse<Network[]>> {
        const header = getHeaders(nodeId);
        return post<Network[]>("/api/docker/networks", null, header);
    },
}


export class DockerUtils {
    static getShortImageId(id: string): string {
        if (id.length < 12) {
            return id;
        }
        if (id.startsWith("sha256:")) {
            return id.substring(7, 19);
        }

        return id.substring(0, 12);
    }

    static getImageRepository(repoDigests?: string[], repoTags?: string[]): { repository: string, tag: string } {
        let repository = '<none>';
        let tag = '<none>';
        // get repoTags first
        if (repoTags && repoTags.length > 0) {
            const [repo, t] = repoTags[0]?.split(':') ?? ['<none>', '<none>'];
            repository = repo ?? '<none>';
            tag = t ?? '<none>';
        }
        // get repoDigests second
        if (repoDigests && repoDigests.length > 0) {
            const [repo, _] = repoDigests[0]?.split('@') ?? ['<none>'];
            repository = repo ?? '<none>';
            tag = '<none>';
        }
        return { repository, tag };
    }

    static formatDisplaySize(size: number): string {
        if (size > 1024 * 1024) {
            return (size / (1024 * 1024)).toFixed(2) + ' MB';
        }
        return (size / 1024).toFixed(2) + ' KB';
    }

    static getShortNetworkId(id: string): string {
        if (id.length < 12) {
            return id;
        }
        return id.substring(0, 12);
    }
}

