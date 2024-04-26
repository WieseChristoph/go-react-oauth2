import { type FC } from "react";
import useAuth from "@/hooks/useAuth";

import UserProfile from "@/components/UserProfile";

const Home: FC = () => {
  const { isAuthenticated, isLoading, user } = useAuth();

  return (
    <div className="flex flex-col items-center gap-2 text-white">
      <h1 className="text-6xl">Hello World</h1>
      {isLoading ? (
        <p>Loading...</p>
      ) : isAuthenticated && user ? (
        <>
          <UserProfile user={user} />
          <a
            href="/auth/logout"
            className="rounded-sm bg-red-900 px-3 py-2 font-bold"
          >
            Logout
          </a>
        </>
      ) : (
        <>
          <a
            href="/auth/discord"
            className="rounded-sm bg-[#5865F2] px-3 py-2 font-bold"
          >
            Login with Discord
          </a>
          <a
            href="/auth/google"
            className="rounded-sm bg-white px-3 py-2 font-bold text-black"
          >
            Login with <span className="text-[#4285F4]">G</span>
            <span className="text-[#EA4335]">o</span>
            <span className="text-[#FBBC04]">o</span>
            <span className="text-[#4285F4]">g</span>
            <span className="text-[#34A853]">l</span>
            <span className="text-[#EA4335]">e</span>
          </a>
        </>
      )}
    </div>
  );
};

export default Home;
