export enum ROLES {
  ADMIN = "ADMIN",
  USER = "USER",
}

export type Role = keyof typeof ROLES;

export type User = {
  objectId: string;
  tenantId: string;
  displayName: string;
  username: string;
  role: Role;
};
