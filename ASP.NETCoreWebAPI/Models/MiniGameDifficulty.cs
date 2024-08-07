namespace ASP.NETCoreWebAPI.Models
{
    public class MiniGameDifficulty
    {
        public int Id { get; set; } // Primary Key
        public required string Name { get; set; } // Key Attribute - Required, names like "I", "II", "III", "IV"
        public GraphicalAsset? Icon { get; set; } // Foreign Key (GraphicalAsset) - Nullable, representing the icon for this difficulty
        public required string OverwritingSettings { get; set; } // Attribute - Required, JSON settings that overwrite the default settings (Remain required or become nullable?)
        public required MiniGame MiniGame { get; set; } // Foreign Key (MiniGame) - Required, links to the parent MiniGame  
    }
}
