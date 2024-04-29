import User from "@/types/User";

const ProfileField: React.FC<{ label: string; value: string }> = ({
  label,
  value,
}) => (
  <p className="flex justify-between gap-5">
    <span className="font-bold">{label}</span>
    <span>{value}</span>
  </p>
);

const UserProfile: React.FC<{ user: User }> = ({ user }) => {
  return (
    <div className="flex items-center gap-3 border border-white p-3 text-white">
      <img
        src={user?.avatar}
        alt="Profile picture"
        referrerPolicy="no-referrer"
      />
      <div>
        <ProfileField label="ID" value={user.id.toString()} />
        <ProfileField label="Name" value={user.name} />
        <ProfileField label="Display Name" value={user.display_name} />
        <ProfileField label="E-Mail" value={user.email} />
        <ProfileField label="Role" value={user.role} />
        <ProfileField label="Created At" value={user.created_at} />
        <ProfileField label="Updated At" value={user.updated_at} />
      </div>
    </div>
  );
};

export default UserProfile;
