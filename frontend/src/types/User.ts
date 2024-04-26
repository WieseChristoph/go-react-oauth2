export enum Role {
  ADMIN = "admin",
  USER = "user",
}

interface User {
  id: number;
  name: string;
  display_name: string;
  email: string;
  avatar: string;
  role: Role;
  created_at: string;
  updated_at: string;
}

export default User;
