export interface DiscordServerImgProps {
  guildID: number | string;
  imgHash: string;
  size?: number;
  className?: string;
}

export default function DiscordServerImg({
  guildID,
  imgHash,
  size,
  className,
}: DiscordServerImgProps) {
  if (size === undefined) {
    size = 128;
  }
  const iconURL = `https://cdn.discordapp.com/icons/${guildID}/${imgHash}.webp?size=${size}`;
  let onHover: React.MouseEventHandler<HTMLImageElement> = () => {};
  let offHover: React.MouseEventHandler<HTMLImageElement> = () => {};
  if (imgHash.startsWith("a_")) {
    const gifURL = iconURL.replace(".webp", ".gif");
    onHover = (e) => {
      e.currentTarget.setAttribute("src", gifURL);
    };
    offHover = (e) => {
      e.currentTarget.setAttribute("src", iconURL);
    };
  }
  return (
    <img
      className={className}
      src={iconURL}
      alt="Discord server icon"
      onMouseOver={onHover}
      onMouseOut={offHover}
    />
  );
}
