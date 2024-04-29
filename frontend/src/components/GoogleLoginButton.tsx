const GoogleLoginButton: React.FC = () => {
  return (
    <a
      className="rounded-sm bg-white px-3 py-2 text-black"
      href="/auth/google"
      aria-label="Login with Google"
    >
      Login with{" "}
      {"Google".split("").map((letter, index) => (
        <span
          key={index}
          className={`font-bold ${index === 0 || index === 3 ? "text-[#4285F4]" : index === 1 || index === 5 ? "text-[#EA4335]" : index === 2 ? "text-[#FBBC04]" : "text-[#34A853]"}`}
        >
          {letter}
        </span>
      ))}
    </a>
  );
};

export default GoogleLoginButton;
