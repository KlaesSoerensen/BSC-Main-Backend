namespace BSC_Main_Backend.Models;

public class AssetUseCase
{
    public static readonly AssetUseCase Environment = new("environment"),
        Player = new("player"), 
        Icon = new("icon"), 
        Unknown = new("unknown");

    public static AssetUseCase ValueOf(string input)
    {
        return input switch
        {
            "environment" => Environment,
            "player" => Player,
            "icon" => Icon,
            _ => Unknown
        };
    }

    public readonly string Value;
    private AssetUseCase(string value)
    {
        Value = value;
    }
}