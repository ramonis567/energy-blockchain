import { Link, useLocation } from "react-router-dom";

const links = [
  { path: "/", label: "Dashboard" },
  { path: "/agents", label: "Agents" },
  { path: "/tokens", label: "Tokens" },
  { path: "/offers", label: "Offers" },
  { path: "/contracts", label: "Contracts" },
];

export default function Sidebar() {
  const location = useLocation();

  return (
    <aside className="w-56 bg-[var(--blueColor)] text-[var(--surfaceColor)] flex flex-col">
      <div className="p-4 font-bold text-2xl">Menu</div>
      <nav className="flex-1">
        {links.map((link) => (
          <Link
            key={link.path}
            to={link.path}
            className={`block px-4 py-2 hover:bg-[var(--highlightColor)] transition-colors ${
              location.pathname === link.path ? "bg-[var(--highlightColor)]" : ""
            }`}
          >
            {link.label}
          </Link>
        ))}
      </nav>
    </aside>
  );
}
