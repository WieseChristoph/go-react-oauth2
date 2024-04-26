import { FC } from "react";
import User from "@/types/User";

const UserProfile: FC<{ user: User }> = ({ user }) => {
  return (
    <div className="flex items-center gap-3 border border-white p-3 text-white">
      <img
        src={user?.avatar}
        alt="Profile picture"
        referrerPolicy="no-referrer"
      />
      <div>
        <p className="flex justify-between gap-5">
          <span className="font-bold">ID</span>
          <span className="">{user.id}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">Name</span>
          <span>{user.name}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">Display Name</span>
          <span>{user.display_name}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">E-Mail</span>
          <span>{user.email}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">Role</span>
          <span>{user.role}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">Created At</span>
          <span>{user.created_at}</span>
        </p>
        <p className="flex justify-between gap-5">
          <span className="font-bold">Updated At</span>
          <span>{user.updated_at}</span>
        </p>
      </div>
    </div>
  );
};

export default UserProfile;
