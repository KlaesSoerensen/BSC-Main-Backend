namespace ASP.NETCoreWebAPI.Models
{
    public class MiniGame
    {
        public int Id { get; set; } // Primary Key
        public required string Name { get; set; } // Key Attribute - Required
        public string? Description { get; set; } // Attribute - Nullable
        public GraphicalAsset? Icon { get; set; } // Foreign Key - Nullable, representing the icon
        public required string Settings { get; set; } // Attribute - Required, JSON settings blob for the minigame (What type of object excatly?)
        public MiniGameDifficulty[]? Difficulties { get; set; } // Multi-Valued Attribute - Nullable, represents various difficulties
    }
}
