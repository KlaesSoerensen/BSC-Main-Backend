namespace ASP.NETCoreWebAPI.Models
{
    public class Location
    {
        public int Id { get; set; } // Primary Key
        public required string Name { get; set; } // Attribute - Required
        public string? Description { get; set; } // Attribute - Nullable
        public int Level { get; set; } // Attribute
        public required AssetCollection AssetCollection { get; set; } // Foreign Key (AssetCollection) - Required
        public MiniGame? Minigame { get; set; } // Foreign Key (Minigame) - Nullable

    }
}
