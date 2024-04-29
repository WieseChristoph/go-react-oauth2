const LogoutButton: React.FC = () => {
  return (
    <a
      href="/auth/logout"
      className="rounded-sm bg-red-900 px-3 py-2 font-bold"
    >
      Logout
    </a>
  );
};

export default LogoutButton;
