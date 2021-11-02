import * as k3s from "@pulumi/k3s";

const random = new k3s.Random("my-random", { length: 24 });

export const output = random.result;