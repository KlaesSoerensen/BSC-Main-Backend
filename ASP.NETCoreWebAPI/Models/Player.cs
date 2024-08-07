namespace ASP.NETCoreWebAPI.Models
{
    public class Player
    {
        public int Id { get; set; } // Primary Key
        public required string IGN { get; set; } // Attribute - Required, In-Game Name of the player
        public GraphicalAsset? Sprite { get; set; } // Foreign Key (GraphicalAsset) - Nullable, link or identifier for the player's sprite
        public Achievement[]? Achievements { get; set; } // Multi-Valued Attribute - Nullable, list of achievements associated with the player
    }
}
